package nbfx

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
)

type record interface {
	isElementStart() bool
	isAttribute() bool
	getName() string
	read(reader *bytes.Reader) (xml.Token, error)
	write(writer io.Writer) error
}

func getRecord(codec *codec, b byte) record {
	if recordFunc, ok := records[b]; ok {
		return recordFunc(codec)
	}

	return nil
}

type startElementRecord struct{}

func (r *startElementRecord) isElementStart() bool { return true }
func (r *startElementRecord) isAttribute() bool    { return false }

type attributeRecord struct{}

func (r *attributeRecord) isElementStart() bool { return false }
func (r *attributeRecord) isAttribute() bool    { return true }

type textRecord struct{}

func (r *textRecord) isElementStart() bool { return false }
func (r *textRecord) isAttribute() bool    { return false }

var records = map[byte]func(*codec) record{
	0x06: func(codec *codec) record { return &shortDictionaryAttributeRecord{codec: codec} },
	0x0B: func(codec *codec) record { return &dictionaryXmlnsAttributeRecord{codec: codec} },
	//0x0C-0x25: func(codec *codec) record { return &prefixDictionaryAttributeAZRecord{codec: codec, prefixIndex: 0x0C-0x25}}, ADDED IN init()
	0x40: func(codec *codec) record { return &shortElementRecord{codec: codec} },
	0x41: func(codec *codec) record { return &elementRecord{codec: codec} },
	0x43: func(codec *codec) record { return &dictionaryElementRecord{codec: codec} },
	//0x44-0x5D: func(codec *codec) record { return &prefixDictionaryElementAZRecord{codec: codec, prefixIndex: 0x44-0x5D}}, ADDED IN init()
	//0x5E-0x77: func(codec *codec) record { return &prefixElementAZRecord{codec: codec, prefixIndex: 0x5E-0x77}}, ADDED IN init()
	0x82: func(codec *codec) record { return &oneTextRecord{codec: codec} },
	0x99: func(codec *codec) record { return &chars8TextWithEndElementRecord{codec: codec} },
}

func init() {
	for b := 0; b < 26; b++ {
		byt := byte(b)
		records[byte(0x0C+byt)] = func(codec *codec) record { return &prefixDictionaryAttributeAZRecord{codec: codec, prefixIndex: byt} }
		records[byte(0x44+byt)] = func(codec *codec) record { return &prefixDictionaryElementAZRecord{codec: codec, prefixIndex: byt} }
		records[byte(0x5E+byt)] = func(codec *codec) record { return &prefixElementAZRecord{codec: codec, prefixIndex: byt} }
	}
}

//(0x06)
type shortDictionaryAttributeRecord struct {
	codec *codec
	*attributeRecord
	nameIndex uint32
}

func (r *shortDictionaryAttributeRecord) getName() string {
	return "ShortDictionaryAttributeRecord (0x06)"
}

func (r *shortDictionaryAttributeRecord) read(reader *bytes.Reader) (xml.Token, error) {
	return nil, errors.New("NotImplemented: shortDictionaryAttributeRecord.write")
}

func (r *shortDictionaryAttributeRecord) write(writer io.Writer) error {
	return errors.New("NotImplemented: shortDictionaryAttributeRecord.write")
}

//0x99
type chars8TextWithEndElementRecord struct {
	codec *codec
	*elementRecord
	name string
}

func (r *chars8TextWithEndElementRecord) getName() string {
	return "chars8TextWithEndElementRecord (0x99)"
}

func (r *chars8TextWithEndElementRecord) read(reader *bytes.Reader) (xml.Token, error) {
	return nil, errors.New("NotImplemented: chars8TextWithEndElementRecord.read")
}

func (r *chars8TextWithEndElementRecord) write(writer io.Writer) error {
	return errors.New("NotImplemented: chars8TextWithEndElementRecord.write")
}

//(0x82)
type oneTextRecord struct {
	*textRecord
	codec *codec
}

func (r *oneTextRecord) getName() string {
	return "OneText (0x82)"
}

func (r *oneTextRecord) read(reader *bytes.Reader) (xml.Token, error) {
	return "1", nil
}

func (r *oneTextRecord) write(writer io.Writer) error {
	_, err := writer.Write([]byte("1"))
	return err
}

//(0x0C-0x25)
type prefixDictionaryAttributeAZRecord struct {
	codec *codec
	*attributeRecord
	prefixIndex byte
	nameIndex   uint32
}

func (r *prefixDictionaryAttributeAZRecord) getName() string {
	return fmt.Sprintf("PrefixDictionaryAttribute%s (%#x)", string(byte('A')+r.prefixIndex), r.prefixIndex)
}

func (r *prefixDictionaryAttributeAZRecord) read(reader *bytes.Reader) (xml.Token, error) {
	name, err := readDictionaryString(reader, r.codec)
	if err != nil {
		return nil, err
	}
	return xml.Attr{Name: xml.Name{Local: string(byte('a'+byte(r.prefixIndex))) + ":" + name}}, nil
}

func (r *prefixDictionaryAttributeAZRecord) write(writer io.Writer) error {
	writer.Write([]byte{0x0C + r.prefixIndex})
	_, err := writeMultiByteInt31(writer, r.nameIndex)
	return err
}

//(0x44-0x5D)
type prefixDictionaryElementAZRecord struct {
	codec *codec
	*startElementRecord
	prefixIndex byte
	nameIndex   uint32
}

func (r *prefixDictionaryElementAZRecord) getName() string {
	return fmt.Sprintf("PrefixDictionaryElement%s (%#x)", string(byte('A')+r.prefixIndex), r.prefixIndex)
}

func (r *prefixDictionaryElementAZRecord) read(reader *bytes.Reader) (xml.Token, error) {
	name, err := readDictionaryString(reader, r.codec)
	if err != nil {
		return nil, err
	}
	return xml.StartElement{Name: xml.Name{Local: string(byte('a'+byte(r.prefixIndex))) + ":" + name}}, nil
}

func (r *prefixDictionaryElementAZRecord) write(writer io.Writer) error {
	writer.Write([]byte{0x44 + r.prefixIndex})
	_, err := writeMultiByteInt31(writer, r.nameIndex)
	return err
}

//(0x43)
type dictionaryElementRecord struct {
	codec *codec
	*startElementRecord
	nameIndex uint32
	prefix    string
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

func (r *dictionaryElementRecord) write(writer io.Writer) error {
	return errors.New("NotImplemented: dictionaryElementRecord.write")
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

	return xml.Attr{Name: xml.Name{Local: "xmlns:" + name}, Value: val}, nil
}

func (r *dictionaryXmlnsAttributeRecord) write(writer io.Writer) error {
	return errors.New("NotImplemented: dictionaryXmlnsAttributeRecord.write")
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

func (r *shortElementRecord) write(writer io.Writer) error {
	return errors.New("NotImplemented: shortElementRecord.write")
}

// 0x5E-0x77
type prefixElementAZRecord struct {
	codec *codec
	*startElementRecord
	name        string
	prefixIndex byte
}

func (r *prefixElementAZRecord) getName() string {
	return "prefixElementAZRecord (0x5E-0x77)"
}

func (r *prefixElementAZRecord) read(reader *bytes.Reader) (xml.Token, error) {
	return nil, errors.New("NotImplemented: prefixElementAZRecord.read")
}

func (r *prefixElementAZRecord) write(writer io.Writer) error {
	return errors.New("NotImplemented: prefixElementAZRecord.write")
}

// 0x41
type elementRecord struct {
	codec *codec
	*startElementRecord
	name   string
	prefix string
}

func (r *elementRecord) getName() string {
	return "elementRecord (0x41)"
}

func (r *elementRecord) read(reader *bytes.Reader) (xml.Token, error) {
	return nil, errors.New("NotImplemented: elementRecord.read")
}

func (r *elementRecord) write(writer io.Writer) error {
	return errors.New("NotImplemented: elementRecord.write")
}
