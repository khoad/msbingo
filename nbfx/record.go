package nbfx

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

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
	prefix string
	name string
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

//(0x43)
type dictionaryElement struct {
	codec *codec
	name string
}

func (r *dictionaryElement) isElementStart() bool{
	return true
}

func (r *dictionaryElement) isAttribute() bool {
	return false
}

func (r *dictionaryElement) getName() string {
	return "DictionaryElement (0x43)"
}

func (r *dictionaryElement) read(reader *bytes.Reader) (xml.Token, error) {
	name, err := readDictionaryString(reader, r.codec)
	if err != nil {
		return nil, err
	}
	return xml.StartElement{Name:xml.Name{Local: name}}, nil
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
