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
		return &prefixDictionaryElementS{codec: codec}
	} else if b == 0x0B {
		return &dictionaryXmlnsAttribute{codec}
	}
	return nil
}

//(0x56)
type prefixDictionaryElementS struct {
	codec       *codec
	prefix      string
	prefixIndex int
	name        string
	nameIndex   uint32
}

func (r *prefixDictionaryElementS) isElementStart() bool {
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
	return xml.StartElement{Name: xml.Name{Local: "s:" + name}}, nil
}

//(0x43)
type dictionaryElement struct {
	codec     *codec
	name      string
	nameIndex uint32
	prefix    string
}

func (r *dictionaryElement) isElementStart() bool {
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
	return xml.StartElement{Name: xml.Name{Local: name}}, nil
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

func (r *dictionaryXmlnsAttribute) isElementStart() bool {
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

	return xml.Attr{Name: xml.Name{Local: "xmlns:" + name}, Value: val}, nil
}

//(0x40)
type shortElement struct {
	codec *codec
	name  string
}

func (r *shortElement) isElementStart() bool {
	return false
}

func (r *shortElement) isAttribute() bool {
	return true
}

func (r *shortElement) getName() string {
	return "shortElement (0x40)"
}

func (r *shortElement) read(reader *bytes.Reader) (xml.Token, error) {
	panic("NIE")
}

//??
type prefixShortElement struct {
	codec  *codec
	prefix string
}

func (r *prefixShortElement) isElementStart() bool {
	return false
}

func (r *prefixShortElement) isAttribute() bool {
	return true
}

func (r *prefixShortElement) getName() string {
	return "prefixShortElement (??)"
}

func (r *prefixShortElement) read(reader *bytes.Reader) (xml.Token, error) {
	panic("NIE")
}
