package nbfx

import (
	"bytes"
	"errors"
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
		return &prefixDictionaryElementSRecord{codec: codec}
	} else if b == 0x0B {
		return &dictionaryXmlnsAttributeRecord{codec: codec}
	}
	return nil
}

type startElementRecord struct {}
func (r *startElementRecord) isElementStart() bool { return true; }
func (r *startElementRecord) isAttribute() bool { return false; }

type attributeRecord struct {}
func (r *attributeRecord) isElementStart() bool { return false }
func (r *attributeRecord) isAttribute() bool { return true }

//(0x56)
type prefixDictionaryElementSRecord struct {
	codec *codec
	*startElementRecord
	prefixIndex int
	nameIndex uint32
}

func (r *prefixDictionaryElementSRecord) getName() string {
	return "PrefixDictionaryElementS (0x56)"
}

func (r *prefixDictionaryElementSRecord) read(reader *bytes.Reader) (xml.Token, error) {
	name, err := readDictionaryString(reader, r.codec)
	if err != nil {
		return nil, err
	}
	return xml.StartElement{Name: xml.Name{Local: "s:" + name}}, nil
}

//(0x43)
type dictionaryElementRecord struct {
	codec *codec
	*startElementRecord
	nameIndex uint32
	prefix string
}

func (r *dictionaryElementRecord) getName() string {
	return "DictionaryElement (0x43)"
}

func (r *dictionaryElementRecord) read(reader *bytes.Reader) (xml.Token, error) {
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
type dictionaryXmlnsAttributeRecord struct {
	codec *codec
	*attributeRecord
}

func (r *dictionaryXmlnsAttributeRecord) getName() string {
	return "dictionaryXmlnsAttribute (0x0B)"
}

func (r *dictionaryXmlnsAttributeRecord) read(reader *bytes.Reader) (xml.Token, error) {
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

// 0x40
type shortElementRecord struct {
	codec *codec
	*startElementRecord
	name string
}

func (r *shortElementRecord) getName() string {
	return "shortElementRecord (0x40)"
}

func (r *shortElementRecord) read(reader *bytes.Reader) (xml.Token, error) {
	return nil, errors.New("NotImplemented: shortElementRecord.read")
}

// 0x5E-0x77
type prefixElementAZRecord struct {
	codec *codec
	*startElementRecord
	name string
	prefixIndex int
}

func (r *prefixElementAZRecord) getName() string {
	return "prefixElementAZRecord (0x5E-0x77)"
}

func (r *prefixElementAZRecord) read(reader *bytes.Reader) (xml.Token, error) {
	return nil, errors.New("NotImplemented: prefixElementAZRecord.read")
}

// 0x41
type elementRecord struct {
	codec *codec
	*startElementRecord
	name string
	prefix string
}

func (r *elementRecord) getName() string {
	return "elementRecord (0x41)"
}

func (r *elementRecord) read(reader *bytes.Reader) (xml.Token, error) {
	return nil, errors.New("NotImplemented: elementRecord.read")
}
