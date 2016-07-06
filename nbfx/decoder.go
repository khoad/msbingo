package nbfx

import (
	"bytes"
	"encoding/xml"
	"io"
	"fmt"
	"encoding/binary"
	"math"
)

type decoder struct {
	codec        codec
	elementStack Stack
}

func NewDecoder() Decoder {
	return NewDecoderWithStrings(nil)
}

func NewDecoderWithStrings(dictionaryStrings map[uint32]string) Decoder {
	decoder := &decoder{codec{make(map[uint32]string), make(map[string]uint32)}, Stack{}}
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
	rec, err := getNextRecord(d, reader)
	for err == nil && rec != nil {
		if rec.isElement() {
			elementReader := rec.(elementRecordDecoder)
			rec, err = elementReader.decodeElement(xmlEncoder, reader)
		} else { // text record
			//fmt.Println("Expecting Text Record and got", rec.getName())
			textReader := rec.(textRecordDecoder)
			_, err = textReader.decodeText(xmlEncoder, reader)
			rec = nil
		}
		if err == nil && rec == nil {
			rec, err = getNextRecord(d, reader)
		}
	}
	//fmt.Println("Exiting main decoder loop...")
	xmlEncoder.Flush()
	if err != nil && err != io.EOF {
		return xmlBuf.String(), err
	}
	return xmlBuf.String(), nil
}

func readMultiByteInt31(reader *bytes.Reader) (uint32, error) {
	b, err := reader.ReadByte()
	if uint32(b) < MASK_MBI31 {
		return uint32(b), err
	}
	nextB, err := readMultiByteInt31(reader)
	return MASK_MBI31*(nextB-1) + uint32(b), err
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
		var err error
		buf, err := readBytes(reader, 1)
		if err != nil {
			return uint32(0), err
		}
		var val uint8
		binary.Read(buf, binary.LittleEndian, &val)
		return uint32(val), err
	})
}

func readChars16Text(reader *bytes.Reader) (string, error) {
	return readStringBytes(reader, func(r *bytes.Reader) (uint32, error) {
		var err error
		buf, err := readBytes(reader, 2)
		if err != nil {
			return uint32(0), err
		}
		var val uint16
		binary.Read(buf, binary.LittleEndian, &val)
		return uint32(val), err
	})
}

func readChars32Text(reader *bytes.Reader) (string, error) {
	return readStringBytes(reader, func(r *bytes.Reader) (uint32, error) {
		var err error
		buf, err := readBytes(reader, 4)
		if err != nil {
			return uint32(0), err
		}
		var val uint32
		binary.Read(buf, binary.LittleEndian, &val)
		return uint32(val), err
	})
}

func readBytes(reader *bytes.Reader, numBytes int) (*bytes.Buffer, error) {
	var err error
	buf := &bytes.Buffer{}
	for i := 0; i < numBytes && err == nil; i++ {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		buf.WriteByte(b)
		if err != nil {
			return nil, err
		}
	}
	return buf, nil
}

func readInt8Text(reader *bytes.Reader) (string, error) {
	var err error
	buf, err := readBytes(reader, 1)
	if err != nil {
		return "", err
	}
	var val int8
	binary.Read(buf, binary.LittleEndian, &val)
	return fmt.Sprintf("%d", val), nil
}

func readInt16Text(reader *bytes.Reader) (string, error) {
	var err error
	buf, err := readBytes(reader, 2)
	if err != nil {
		return "", err
	}
	var val int16
	binary.Read(buf, binary.LittleEndian, &val)
	return fmt.Sprintf("%d", val), nil
}

func readInt32Text(reader *bytes.Reader) (string, error) {
	var err error
	buf, err := readBytes(reader, 4)
	if err != nil {
		return "", err
	}
	var val int32
	binary.Read(buf, binary.LittleEndian, &val)
	return fmt.Sprintf("%d", val), nil
}

func readInt64Text(reader *bytes.Reader) (string, error) {
	var err error
	buf, err := readBytes(reader, 8)
	if err != nil {
		return "", err
	}
	var val int64
	binary.Read(buf, binary.LittleEndian, &val)
	return fmt.Sprintf("%d", val), nil
}

func readFloatText(reader *bytes.Reader) (string, error) {
	var err error
	buf, err := readBytes(reader, 4)
	if err != nil {
		return "", err
	}
	var val float32
	binary.Read(buf, binary.LittleEndian, &val)
	if val == float32(math.Inf(1)) {
		return "INF", nil
	} else if val == float32(math.Inf(-1)) {
		return "-INF", nil
	}
	return fmt.Sprintf("%v", val), nil
}

func readDoubleText(reader *bytes.Reader) (string, error) {
	var err error
	buf, err := readBytes(reader, 8)
	if err != nil {
		return "", err
	}
	var val float64
	binary.Read(buf, binary.LittleEndian, &val)
	if val == math.Inf(1) {
		return "INF", nil
	} else if val == math.Inf(-1) {
		return "-INF", nil
	}
	return fmt.Sprintf("%v", val), nil
}
