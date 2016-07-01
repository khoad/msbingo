package nbfx

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
)

type record interface {
	isElement() bool
	isAttribute() bool
	getName() string
}

type elementRecordReader interface {
	record
	readElement(x xml.Encoder, reader *bytes.Reader) (record, error)
}

type elementRecordWriter interface {
	record
	writeElement(writer io.Writer) error
}

type attributeRecordReader interface {
	record
	readAttribute(x xml.Encoder, reader *bytes.Reader) (xml.Attr, bool, error)
}

type textRecordReader interface {
	record
	readText(reader *bytes.Reader) (string, bool, error)
}

func readRecord(decoder *decoder, reader *bytes.Reader) (record, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	if recordFunc, ok := records[b]; ok {
		record := recordFunc(decoder)
		//fmt.Println("record.getName()", record.getName())
		return record, nil
	}

	return nil, errors.New(fmt.Sprintf("Unknown record %#x", b))
}

type elementRecordBase struct {
	decoder *decoder
}

func (r *elementRecordBase) isElement() bool   { return true }
func (r *elementRecordBase) isAttribute() bool { return false }

func (r *elementRecordBase) readElementAttributes(element xml.StartElement, x xml.Encoder, reader *bytes.Reader) (record, error) {
	// get next record
	//fmt.Printf("getting next record")
	rec, err := readRecord(r.decoder, reader)

	var peekRecord record

	var attributeToken xml.Attr
	closeElement := false
	for rec != nil && !closeElement {
		if err != nil {
			//fmt.Println("Error getting next record", err.Error())
			return nil, err
		}

		var attrReader attributeRecordReader
		if rec.isAttribute() {
			//fmt.Print("record is attribute", record)
			attrReader = rec.(attributeRecordReader)

			attributeToken, closeElement, err = attrReader.readAttribute(x, reader)
			if err != nil {
				return nil, err
			}
			element.Attr = append(element.Attr, attributeToken)

			rec, err = readRecord(r.decoder, reader)
		} else {
			//fmt.Print("record is NOT attribute", record)
			attrReader = nil
			peekRecord = rec
			rec = nil
		}
	}

	//fmt.Printf("Encoding element %s", element)
	err = x.EncodeToken(element)
	if err != nil {
		return nil, err
	}

	r.decoder.elementStack.Push(element)
	fmt.Println("Pushed ", element.Name)

	return peekRecord, nil
}

type attributeRecordBase struct{
	decoder *decoder
}

func (r *attributeRecordBase) isElement() bool   { return false }
func (r *attributeRecordBase) isAttribute() bool { return true }

type textRecordBase struct {
	decoder          *decoder
	withEndElement bool
	textName       string
	recordId       byte
	charData       func(*bytes.Reader) (string, error)
}

func (r *textRecordBase) isElement() bool   { return false }
func (r *textRecordBase) isAttribute() bool { return false }

func (r *textRecordBase) getName() string {
	return fmt.Sprintf("%s (%#x)", r.textName, r.recordId)
}

func (r *textRecordBase) readText(reader *bytes.Reader) (string, bool, error) {
	text, err := r.charData(reader)
	if err != nil {
		return "", false, err
	}
	return text, false, nil
}

var records = map[byte]func(*decoder) record{
	0x01: func(decoder *decoder) record { return &endElementRecord{&elementRecordBase{decoder: decoder}} },
	0x06: func(decoder *decoder) record { return &shortDictionaryAttributeRecord{&attributeRecordBase{decoder: decoder}, 0} },
	0x0B: func(decoder *decoder) record { return &dictionaryXmlnsAttributeRecord{&attributeRecordBase{decoder: decoder}} },
	//0x0C-0x25: func(decoder *decoder) record { return &prefixDictionaryAttributeAZRecord{decoder: decoder, prefixIndex: 0x0C-0x25}}, ADDED IN init()
	0x40: func(decoder *decoder) record { return &shortElementRecord{&elementRecordBase{decoder: decoder}, ""} },
	0x41: func(decoder *decoder) record { return &elementRecord{&elementRecordBase{decoder: decoder}, "", ""} },
	0x43: func(decoder *decoder) record { return &dictionaryElementRecord{&elementRecordBase{decoder: decoder}, 0, ""} },
	//0x44-0x5D: func(decoder *decoder) record { return &prefixDictionaryElementAZRecord{decoder: decoder, prefixIndex: 0x44-0x5D}}, ADDED IN init()
	//0x5E-0x77: func(decoder *decoder) record { return &prefixElementAZRecord{decoder: decoder, prefixIndex: 0x5E-0x77}}, ADDED IN init()
	//0x80: func(decoder *decoder) record {
	//	return &textRecord{decoder: decoder, withEndElement: false, textName: "ZeroText", recordId: 0x80, charData: "0"}
	//},
	//0x81: func(decoder *decoder) record { return &zeroTextRecord{decoder: decoder, withEndElement: true} },
	//0x82: func(decoder *decoder) record { return &oneTextRecord{decoder: decoder} },
	//0x99: func(decoder *decoder) record { return &chars8TextWithEndElementRecord{decoder: decoder} },
}

func addTextRecord(recordId byte, textName string, charData func(*bytes.Reader) (string, error)) {
	records[recordId] = func(decoder *decoder) record {
		return &textRecordBase{decoder: decoder, withEndElement: false, textName: textName, recordId: recordId, charData: charData}
	}
	records[recordId+1] = func(decoder *decoder) record {
		return &textRecordBase{decoder: decoder, withEndElement: true, textName: textName + "WithEndElement", recordId: recordId + 1, charData: charData}
	}
}

func init() {
	for b := 0; b < 26; b++ {
		byt := byte(b)
		records[byte(0x0C+byt)] = func(decoder *decoder) record { return &prefixDictionaryAttributeAZRecord{&attributeRecordBase{decoder: decoder}, byt, 0} }
		records[byte(0x44+byt)] = func(decoder *decoder) record {
			return &prefixDictionaryElementAZRecord{&elementRecordBase{decoder: decoder}, byt, 0}
		}
		records[byte(0x5E+byt)] = func(decoder *decoder) record { return &prefixElementAZRecord{&elementRecordBase{decoder: decoder}, "", byt} }
	}
	addTextRecord(0x80, "ZeroText", func(reader *bytes.Reader) (string, error) { return "0", nil })
	addTextRecord(0x82, "OneText", func(reader *bytes.Reader) (string, error) { return "1", nil })
	addTextRecord(0x98, "Chars8Text", func(reader *bytes.Reader) (string, error) { return readChars8Text(reader) })
}

//(0x01)
type endElementRecord struct {
	*elementRecordBase
}

func (r *endElementRecord) getName() string {
	return "EndElementRecord (0x01)"
}

func (r *endElementRecord) readElement(x xml.Encoder, reader *bytes.Reader) (record, error) {
	item := r.decoder.elementStack.Pop()
	element := item.(xml.StartElement)
	fmt.Println("Popped ", element.Name)
	err := x.EncodeToken(xml.EndElement{Name:element.Name})
	return nil, err
}

func (r *endElementRecord) write(writer io.Writer) error {
	return errors.New("NotImplemented: endElementRecord.write")
}

//(0x06)
type shortDictionaryAttributeRecord struct {
	*attributeRecordBase
	nameIndex uint32
}

func (r *shortDictionaryAttributeRecord) getName() string {
	return "ShortDictionaryAttributeRecord (0x06)"
}

func (r *shortDictionaryAttributeRecord) readAttribute(x xml.Encoder, reader *bytes.Reader) (xml.Attr, bool, error) {
	name, err := readDictionaryString(reader, r.decoder)
	if err != nil {
		return xml.Attr{}, false, err
	}
	record, err := readRecord(r.decoder, reader)
	if err != nil {
		return xml.Attr{}, false, err
	}
	textReader := record.(textRecordReader)
	text, closeElement, err := textReader.readText(reader)
	if err != nil {
		return xml.Attr{}, false, err
	}
	return xml.Attr{Name: xml.Name{Local: name}, Value: text}, closeElement, nil
}

////0x99
//type chars8TextWithEndElementRecord struct {
//	decoder *decoder
//	*textRecord
//	name string
//}
//
//func (r *chars8TextWithEndElementRecord) getName() string {
//	return "chars8TextWithEndElementRecord (0x99)"
//}
//
//func (r *chars8TextWithEndElementRecord) read(x xml.Encoder, reader *bytes.Reader) (record, error) {
//	text, err := readString(reader)
//	if err != nil {
//		return nil, err
//	}
//	err = x.EncodeToken(xml.CharData(text))
//	if err != nil {
//		return nil, err
//	}
//	err = x.EncodeToken(xml.EndElement{})
//	return nil, err
//}
//
//func (r *chars8TextWithEndElementRecord) write(writer io.Writer) error {
//	return errors.New("NotImplemented: chars8TextWithEndElementRecord.write")
//}

//
////(0x80-81)
//type zeroTextRecord struct {
//	decoder *decoder
//	*textRecord
//}
//
//func (r *zeroTextRecord) getName() string {
//	return "ZeroText (0x80-81)"
//}
//
//func (r *zeroTextRecord) read(x xml.Encoder, reader *bytes.Reader) (record, error) {
//	err := x.EncodeToken(xml.CharData("0"))
//	if err != nil {
//		return nil, err
//	}
//	if r.withEndElement {
//		err = x.EncodeToken(xml.EndElement{})
//	}
//	return nil, err
//}
//
//func (r *zeroTextRecord) write(writer io.Writer) error {
//	_, err := writer.Write([]byte("0"))
//	return err
//}

////(0x82-83)
//type oneTextRecord struct {
//	*textRecord
//	decoder *decoder
//}
//
//func (r *oneTextRecord) getName() string {
//	return "OneText (0x82)"
//}
//
//func (r *oneTextRecord) read(x xml.Encoder, reader *bytes.Reader) (record, error) {
//	err := x.EncodeToken(xml.CharData("1"))
//
//	if r.withEndElement {
//
//	}
//	return nil, err
//}
//
//func (r *oneTextRecord) write(writer io.Writer) error {
//	_, err := writer.Write([]byte("1"))
//	return err
//}

//(0x0C-0x25)
type prefixDictionaryAttributeAZRecord struct {
	*attributeRecordBase
	prefixIndex byte
	nameIndex   uint32
}

func (r *prefixDictionaryAttributeAZRecord) getName() string {
	return fmt.Sprintf("PrefixDictionaryAttributeAZRecord (%#x)", byte(0x0C+r.prefixIndex))
}

func (r *prefixDictionaryAttributeAZRecord) readAttribute(x xml.Encoder, reader *bytes.Reader) (xml.Attr, bool, error) {
	name, err := readDictionaryString(reader, r.decoder)
	if err != nil {
		return xml.Attr{}, false, err
	}
	attrToken := xml.Attr{Name: xml.Name{Local: string(byte('a'+byte(r.prefixIndex))) + ":" + name}}
	record, err := readRecord(r.decoder, reader)
	if err != nil {
		return xml.Attr{}, false, err
	}
	textRecord := record.(textRecordReader)
	if textRecord == nil {
		return xml.Attr{}, false, errors.New("Expected TextRecord")
	}
	text, closeElement, err := textRecord.readText(reader)
	if err != nil {
		return xml.Attr{}, false, err
	}
	attrToken.Value = text
	return attrToken, closeElement, nil
}

//
//func (r *prefixDictionaryAttributeAZRecord) write(writer io.Writer) error {
//	writer.Write([]byte{0x0C + r.prefixIndex})
//	_, err := writeMultiByteInt31(writer, r.nameIndex)
//	return err
//}

//(0x44-0x5D)
type prefixDictionaryElementAZRecord struct {
	*elementRecordBase
	prefixIndex byte
	nameIndex   uint32
}

func (r *prefixDictionaryElementAZRecord) getName() string {
	return fmt.Sprintf("PrefixDictionaryElement%s (%#x)", string(byte('A')+r.prefixIndex), 0x44+r.prefixIndex)
}

func (r *prefixDictionaryElementAZRecord) readElement(x xml.Encoder, reader *bytes.Reader) (record, error) {
	name, err := readDictionaryString(reader, r.decoder)
	if err != nil {
		return nil, err
	}
	element := xml.StartElement{Name: xml.Name{Local: string(byte('a'+byte(r.prefixIndex))) + ":" + name}}

	return r.readElementAttributes(element, x, reader)
}

func (r *prefixDictionaryElementAZRecord) write(writer io.Writer) error {
	writer.Write([]byte{0x44 + r.prefixIndex})
	_, err := writeMultiByteInt31(writer, r.nameIndex)
	return err
}

//(0x43)
type dictionaryElementRecord struct {
	*elementRecordBase
	nameIndex uint32
	prefix    string
}

func (r *dictionaryElementRecord) getName() string {
	return "DictionaryElement (0x43)"
}

func (r *dictionaryElementRecord) readElement(x xml.Encoder, reader *bytes.Reader) (record, error) {
	name, err := readDictionaryString(reader, r.decoder)
	if err != nil {
		return nil, err
	}
	element := xml.StartElement{Name: xml.Name{Local: name}}

	return r.readElementAttributes(element, x, reader)
}

func readDictionaryString(reader *bytes.Reader, decoder *decoder) (string, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return "", err
	}
	key := uint32(b)
	if val, ok := decoder.codec.dict[key]; ok {
		return val, nil
	}
	return fmt.Sprintf("str%d", b), nil
}

//(0x0B)
type dictionaryXmlnsAttributeRecord struct {
	*attributeRecordBase
}

func (r *dictionaryXmlnsAttributeRecord) getName() string {
	return "dictionaryXmlnsAttribute (0x0B)"
}

func (r *dictionaryXmlnsAttributeRecord) readAttribute(x xml.Encoder, reader *bytes.Reader) (xml.Attr, bool, error) {
	name, err := readString(reader)
	if err != nil {
		return xml.Attr{}, false, err
	}

	val, err := readDictionaryString(reader, r.decoder)
	if err != nil {
		return xml.Attr{}, false, err
	}

	return xml.Attr{Name: xml.Name{Local: "xmlns:" + name}, Value: val}, false, nil
}

func (r *dictionaryXmlnsAttributeRecord) write(writer io.Writer) error {
	return errors.New("NotImplemented: dictionaryXmlnsAttributeRecord.write")
}

// 0x40
type shortElementRecord struct {
	*elementRecordBase
	name string
}

func (r *shortElementRecord) getName() string {
	return "shortElementRecord (0x40)"
}

func (r *shortElementRecord) readElement(x xml.Encoder, reader *bytes.Reader) (record, error) {
	name, err := readString(reader)
	if err != nil {
		return nil, err
	}
	element := xml.StartElement{Name: xml.Name{Local: name}}

	return r.readElementAttributes(element, x, reader)
}

//func (r *shortElementRecord) write(writer io.Writer) error {
//	return errors.New("NotImplemented: shortElementRecord.write")
//}

// 0x5E-0x77
type prefixElementAZRecord struct {
	*elementRecordBase
	name        string
	prefixIndex byte
}

func (r *prefixElementAZRecord) getName() string {
	return fmt.Sprintf("PrefixElementAZRecord (%#x)", r.prefixIndex+0x5E)
}

func (r *prefixElementAZRecord) readElement(x xml.Encoder, reader *bytes.Reader) (record, error) {
	name, err := readString(reader)
	if err != nil {
		return nil, err
	}
	element := xml.StartElement{Name: xml.Name{Local: string(byte('a'+byte(r.prefixIndex))) + ":" + name}}

	return r.readElementAttributes(element, x, reader)
}

//func (r *prefixElementAZRecord) write(writer io.Writer) error {
//	return errors.New("NotImplemented: prefixElementAZRecord.write")
//}

// 0x41
type elementRecord struct {
	*elementRecordBase
	name   string
	prefix string
}

func (r *elementRecord) getName() string {
	return "elementRecord (0x41)"
}

func (r *elementRecord) readElement(x xml.Encoder, reader *bytes.Reader) (record, error) {
	prefix, err := readString(reader)
	if err != nil {
		return nil, err
	}
	name, err := readString(reader)
	if err != nil {
		return nil, err
	}
	element := xml.StartElement{Name: xml.Name{Local: prefix + ":" + name}}

	return r.readElementAttributes(element, x, reader)
}

//func (r *elementRecord) write(writer io.Writer) error {
//	return errors.New("NotImplemented: elementRecord.write")
//}
