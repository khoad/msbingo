package nbfx

import (
	"io"
	"fmt"
	"errors"
	"bytes"
	"encoding/xml"
)

type decoder struct {
	codec codec
}

func NewDecoder() Decoder {
	return NewDecoderWithStrings(nil)
}

func NewDecoderWithStrings(dictionaryStrings map[uint32]string) Decoder {
	decoder := &decoder{codec{make(map[uint32]string)}}
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

type record interface {
	isElementStart() bool
	isAttribute() bool
	getName() string
	read(reader *bytes.Reader) (xml.Token, error)
}

func getRecord(codec *codec, b byte) record {
	if b == 0x56 {
		return &prefixDictionaryElementS{codec}
	} else if b == 0x0B {
		return &dictionaryXmlnsAttribute{codec}
	}
	return nil
}

//(0x56)
type prefixDictionaryElementS struct {
	codec *codec
}

func (r *prefixDictionaryElementS) isElementStart() bool{
	return true
}

func (r *prefixDictionaryElementS) isAttribute() bool {
	return false
}

func (r *prefixDictionaryElementS) getName() string {
	return "PrefixDictionaryElementS (0x56)"
}

func (r *prefixDictionaryElementS) read(reader *bytes.Reader) (xml.Token, error) {
	name, err := readDictionaryString(reader, r.codec)
	if err != nil {
		return nil, err
	}
	return xml.StartElement{Name:xml.Name{Local:"s:" + name}}, nil
}

func readDictionaryString(reader *bytes.Reader, codec *codec) (string, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return "", err
	}
	key := uint32(b)
	if val, ok := codec.dict[key]; ok {
		return val, nil
	}
	return fmt.Sprintf("str%d", b), nil
}

//(0x0B)
type dictionaryXmlnsAttribute struct {
	codec *codec
}

func (r *dictionaryXmlnsAttribute) isElementStart() bool{
	return false
}

func (r *dictionaryXmlnsAttribute) isAttribute() bool {
	return true
}

func (r *dictionaryXmlnsAttribute) getName() string {
	return "dictionaryXmlnsAttribute (0x0B)"
}

func (r *dictionaryXmlnsAttribute) read(reader *bytes.Reader) (xml.Token, error) {
	name, err := readString(reader)
	if err != nil {
		return name, err
	}

	val, err := readDictionaryString(reader, r.codec)
	if err != nil {
		return nil, err
	}
	fmt.Println("Attr", name, val)

	return xml.Attr{Name:xml.Name{Local:"xmlns:" + name}, Value:val}, nil
}

func readMultiByteInt31(reader * bytes.Reader) (uint32, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	return uint32(b), nil //TODO: Handle multibyte values!!!
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
