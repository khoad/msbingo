package nbfx

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
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
	initStartElement := func(token xml.Token) {
		flushStartElement()
		haveStartingElement = true
		startingElement = token.(xml.StartElement)
	}
	for err == nil {
		record := getRecord(&d.codec, b)
		if record == nil {
			xmlEncoder.Flush()
			return xmlBuf.String(), errors.New(fmt.Sprintf("Unknown Record ID %#x", b))
		}
		var token xml.Token
		token, err = record.read(reader)
		if err != nil {
			xmlEncoder.Flush()
			return xmlBuf.String(), err
		}
		if record.isElementStart() {
			initStartElement(token)
		} else if record.isAttribute() {
			startingElement.Attr = append(startingElement.Attr, token.(xml.Attr))
		} else {
			flushStartElement()
			xmlEncoder.EncodeToken(token)
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
	b, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	return uint32(b), nil //TODO: Handle multibyte values!!!
}

func writeMultiByteInt31(writer io.Writer, i uint32) (int, error) {
	b, err := writer.Write([]byte{byte(i)}) //TODO: Handle multibyte values!!!
	if err != nil {
		return b, err
	}
	return b, nil
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
