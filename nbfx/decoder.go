package nbfx

import (
	"bytes"
	"encoding/xml"
	"errors"
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
	b, err := reader.ReadByte()
	var startingElement xml.StartElement
	haveStartingElement := false
	flushStartElement := func() {
		if haveStartingElement {
			xmlEncoder.EncodeToken(startingElement)
		}
		haveStartingElement = false
		startingElement = xml.StartElement{}
	}
	initStartElement := func(tokens []xml.Token) {
		flushStartElement()
		haveStartingElement = true
		startingElement = tokens[0].(xml.StartElement)
	}
	for err == nil {
		record := getRecord(&d.codec, b)
		if record == nil {
			xmlEncoder.Flush()
			return xmlBuf.String(), errors.New(fmt.Sprintf("Unknown Record ID %#x", b))
		}
		fmt.Println(record.getName())
		var tokens []xml.Token
		tokens, err = record.read(reader)
		if err != nil {
			xmlEncoder.Flush()
			return xmlBuf.String(), err
		}
		if record.isElementStart() {
			initStartElement(tokens)
		} else if record.isAttribute() {
			for i, t := range tokens {
				fmt.Println(t)
				if i == 0 {
					startingElement.Attr = append(startingElement.Attr, t.(xml.Attr))
				} else {
					xmlEncoder.EncodeToken(t)
				}
			}
		} else {
			flushStartElement()
			for _, t := range tokens {
				xmlEncoder.EncodeToken(t)
			}
		}

		b, err = reader.ReadByte()
	}
	flushStartElement()
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

func readString(reader *bytes.Reader) (string, error) {
	var len uint32
	len, err := readMultiByteInt31(reader)
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
