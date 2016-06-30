package nbfx

import (
	"bytes"
	"encoding/xml"
	//"errors"
	"fmt"
	"io"
	"math"
)

type decoder struct {
	codec codec
}

func NewDecoder() Decoder {
	return NewDecoderWithStrings(nil)
}

func NewDecoderWithStrings(dictionaryStrings map[uint32]string) Decoder {
	decoder := &decoder{codec{make(map[uint32]string), make(map[string]uint32)}}
	if dictionaryStrings != nil {
		for k, v := range dictionaryStrings {
			decoder.codec.addDictionaryString(k, v)
		}
	}
	return decoder
}

func (d *decoder) Decode(bin []byte) (string, error) {
	reader := bytes.NewReader(bin)
	xmlBuf := &bytes.Buffer{}
	xmlEncoder := xml.NewEncoder(xmlBuf)
	record, err := readRecord(&d.codec, reader)
	for err == nil && record != nil {
		if record.isElement() {
			fmt.Println("Processing element" + record.getName())
			elementReader := record.(elementRecordReader)
			record, err = elementReader.readElement(*xmlEncoder, reader)
		} else if record.isAttribute() {
			attributeReader := record.(attributeRecordReader)
			var attr xml.Attr
			attr, _, err = attributeReader.readAttribute(*xmlEncoder, reader)
			xmlEncoder.EncodeToken(attr)
			record = nil
		} else { // text record
			textReader := record.(textRecordReader)
			var text string
			text, _, err = textReader.readText(*xmlEncoder, reader)
			xmlEncoder.EncodeToken(xml.CharData(text))
			record = nil
		}
		if err == nil && record == nil {
			record, err = readRecord(&d.codec, reader)
		}
	}
	xmlEncoder.Flush()
	if err != nil && err != io.EOF {
		return xmlBuf.String(), err
	}
	return xmlBuf.String(), nil
}

func readMultiByteInt31(reader *bytes.Reader) (uint32, error) {
	buf := new([5]byte)
	keepReading := true
	i := -1
	for keepReading {
		i++
		b, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}
		if b >= MASK_MBI31 {
			b -= MASK_MBI31
			keepReading = true
		} else {
			keepReading = false
		}
		buf[i] = b
	}
	var val uint32
	for ; i >= 0; i-- {
		val += uint32(buf[i]) * uint32(math.Pow(128, float64(i)))
	}
	return val, nil
}

func readStringBytes(reader *bytes.Reader, readLenFunc func(r *bytes.Reader) (uint32, error)) (string, error) {
	len, err := readLenFunc(reader)
	if err != nil {
		return "", err
	}
	strBuffer := bytes.Buffer{}
	for i := uint32(0); i < len; {
		b, err := reader.ReadByte()
		if err != nil {
			return strBuffer.String(), err
		}
		strBuffer.WriteByte(b)
		i++
	}
	return strBuffer.String(), nil
}

func readString(reader *bytes.Reader) (string, error) {
	return readStringBytes(reader, func(r *bytes.Reader) (uint32, error) {
		return readMultiByteInt31(r)
	})
}

func readChars8Text(reader *bytes.Reader) (string, error) {
	return readStringBytes(reader, func(r *bytes.Reader) (uint32, error) {
		len, err := reader.ReadByte()
		return uint32(len), err
	})
}
