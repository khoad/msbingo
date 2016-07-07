package nbfx

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	//"time"
)

type record interface {
	isElement() bool
	isAttribute() bool
	getName() string
}

type elementRecordDecoder interface {
	record
	decodeElement(x *xml.Encoder, reader *bytes.Reader) (record, error)
}

type elementRecordWriter interface {
	record
	writeElement(writer io.Writer, x xml.Token) error
}

type attributeRecordDecoder interface {
	record
	decodeAttribute(x *xml.Encoder, reader *bytes.Reader) (xml.Attr, error)
}

type textRecordDecoder interface {
	record
	decodeText(x *xml.Encoder, reader *bytes.Reader) (string, error)
	readText(reader *bytes.Reader) (string, error)
}

type textRecordWriter interface {
	record
	writeText(w io.Writer) error
}

func getNextRecord(decoder *decoder, reader *bytes.Reader) (record, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	if recordFunc, ok := records[b]; ok {
		rec := recordFunc(decoder)
		//fmt.Println("record.getName()", record.getName())
		return rec, nil
	}

	return nil, errors.New(fmt.Sprintf("Unknown record %#x", b))
}

type elementRecordBase struct {
	decoder *decoder
}

func (r *elementRecordBase) isElement() bool   { return true }
func (r *elementRecordBase) isAttribute() bool { return false }

func (r *elementRecordBase) readElementAttributes(element xml.StartElement, x *xml.Encoder, reader *bytes.Reader) (record, error) {
	var peekRecord record
	var attributeToken xml.Attr

	// get next record
	//fmt.Println("getting next record")
	rec, err := getNextRecord(r.decoder, reader)
	for err == nil && rec != nil {
		//fmt.Println("Processing record", rec.getName())
		if err != nil {
			return nil, err
		}

		var attrReader attributeRecordDecoder
		if rec.isAttribute() {
			attrReader = rec.(attributeRecordDecoder)

			attributeToken, err = attrReader.decodeAttribute(x, reader)
			if err != nil {
				return nil, err
			}
			element.Attr = append(element.Attr, attributeToken)

			rec, err = getNextRecord(r.decoder, reader)
			if err != nil {
				return nil, err
			}
		} else {
			attrReader = nil
			peekRecord = rec
			rec = nil
		}
	}
	//fmt.Println("got next record", peekRecord, err)

	err = x.EncodeToken(element)
	if err != nil {
		return nil, err
	}

	r.decoder.elementStack.Push(element)

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
	charData       func(*bytes.Reader, *decoder) (string, error)
}

func (r *textRecordBase) isElement() bool   { return false }
func (r *textRecordBase) isAttribute() bool { return false }

func (r *textRecordBase) getName() string {
	return fmt.Sprintf("%s (%#x)", r.textName, r.recordId)
}

func (r *textRecordBase) readText(reader *bytes.Reader) (string, error) {
	text, err := r.charData(reader, r.decoder)
	if err != nil {
		return "", err
	}
	return text, nil
}

func (r *textRecordBase) decodeText(x *xml.Encoder, reader *bytes.Reader) (string, error) {
	text, err := r.readText(reader)
	if err != nil {
		return "", err
	}
	charData := xml.CharData([]byte(text))
	x.EncodeToken(charData)
	if r.withEndElement {
		rec, err := getNextRecord(r.decoder, bytes.NewReader([]byte{0x01}))
		endElementReader := rec.(elementRecordDecoder)
		if err == nil {
			endElementReader.decodeElement(x, nil)
		}
	}
	return text, nil
}

func (r *textRecordBase) writeText(w *io.Writer) error {
	return errors.New("NotImplement: writeText for " + r.getName())
}

var records = map[byte]func(*decoder) record{}

func initRecords() {
	// Miscellaneous Records
	records[0x01] = func(decoder *decoder) record { return &endElementRecord{&elementRecordBase{decoder: decoder}} }
	records[0x02] = func(decoder *decoder) record { return &commentRecord{&textRecordBase{decoder: decoder}, ""} }
	records[0x03] = func(decoder *decoder) record { return &arrayRecord{&elementRecordBase{decoder: decoder}} }

	// Attribute Records
	records[0x04] = func(decoder *decoder) record { return &shortAttributeRecord{&attributeRecordBase{decoder: decoder}} }
	records[0x05] = func(decoder *decoder) record { return &attributeRecord{&attributeRecordBase{decoder: decoder}} }
	records[0x06] = func(decoder *decoder) record { return &shortDictionaryAttributeRecord{&attributeRecordBase{decoder: decoder}, 0} }
	records[0x07] = func(decoder *decoder) record { return &dictionaryAttributeRecord{&attributeRecordBase{decoder: decoder}} }
	records[0x08] = func(decoder *decoder) record { return &shortXmlnsAttributeRecord{&attributeRecordBase{decoder: decoder}} }
	records[0x09] = func(decoder *decoder) record { return &xmlnsAttributeRecord{&attributeRecordBase{decoder: decoder}} }
	records[0x0A] = func(decoder *decoder) record { return &shortDictionaryXmlnsAttributeRecord{&attributeRecordBase{decoder: decoder}} }
	records[0x0B] = func(decoder *decoder) record { return &dictionaryXmlnsAttributeRecord{&attributeRecordBase{decoder: decoder}} }
	//0x0C-0x25: func(decoder *decoder) record { return &prefixDictionaryAttributeAZRecord{decoder: decoder, prefixIndex: 0x0C-0x25}}, ADDED IN addAzRecords()
	//0x26-0x3F: func(decoder *decoder) record { return &prefixAttributeAZRecord{decoder: decoder, prefixIndex: 0x2626-0x3F}}, ADDED IN addAzRecords()

	// Element Records
	records[0x40] = func(decoder *decoder) record { return &shortElementRecord{&elementRecordBase{decoder: decoder}, ""} }
	records[0x41] = func(decoder *decoder) record { return &elementRecord{&elementRecordBase{decoder: decoder}, "", ""} }
	records[0x42] = func(decoder *decoder) record { return &shortDictionaryElementRecord{&elementRecordBase{decoder: decoder}} }
	records[0x43] = func(decoder *decoder) record { return &dictionaryElementRecord{&elementRecordBase{decoder: decoder}, 0, ""} }
	//0x44-0x5D: func(decoder *decoder) record { return &prefixDictionaryElementAZRecord{decoder: decoder, prefixIndex: 0x44-0x5D}}, ADDED IN addAzRecords()
	//0x5E-0x77: func(decoder *decoder) record { return &prefixElementAZRecord{decoder: decoder, prefixIndex: 0x5E-0x77}}, ADDED IN addAzRecords()

	// Text Records
	addTextRecord(0x80, "ZeroText", func(r *bytes.Reader, d *decoder) (string, error) { return "0", nil })
	addTextRecord(0x82, "OneText", func(r *bytes.Reader, d *decoder) (string, error) { return "1", nil })
	addTextRecord(0x84, "FalseText", func(r *bytes.Reader, d *decoder) (string, error) { return "false", nil })
	addTextRecord(0x86, "TrueText", func(r *bytes.Reader, d *decoder) (string, error) { return "true", nil})
	addTextRecord(0x88, "Int8Text", func(r *bytes.Reader, d *decoder) (string, error) { return readInt8Text(r) })
	addTextRecord(0x8A, "Int16Text", func(r *bytes.Reader, d *decoder) (string, error) { return readInt16Text(r) })
	addTextRecord(0x8C, "Int32Text", func(r *bytes.Reader, d *decoder) (string, error) { return readInt32Text(r) })
	addTextRecord(0x8E, "Int64Text", func(r *bytes.Reader, d *decoder) (string, error) { return readInt64Text(r) })
	addTextRecord(0x90, "FloatText", func(r *bytes.Reader, d *decoder) (string, error) { return readFloatText(r) })
	addTextRecord(0x92, "DoubleText", func(r *bytes.Reader, d *decoder) (string, error) { return readDoubleText(r) })
	addTextRecord(0x94, "DecimalText", func(r *bytes.Reader, d *decoder) (string, error) { return readDecimalText(r) })
	addTextRecord(0x96, "DateTimText", func(r *bytes.Reader, d *decoder) (string, error) { return readDateTimeText(r) })
	addTextRecord(0x98, "Chars8Text", func(r *bytes.Reader, d *decoder) (string, error) { return readChars8Text(r) })
	addTextRecord(0x9A, "Chars16Text", func(r *bytes.Reader, d *decoder) (string, error) { return readChars16Text(r) })
	addTextRecord(0x9C, "Chars32Text", func(r *bytes.Reader, d *decoder) (string, error) { return readChars32Text(r) })
	addTextRecord(0x9E, "Bytes8Text", func(r *bytes.Reader, d *decoder) (string, error) { return readBytes8Text(r) })
	addTextRecord(0xA0, "Bytes16Text", func(r *bytes.Reader, d *decoder) (string, error) { return readBytes16Text(r) })
	addTextRecord(0xA2, "Bytes32Text", func(r *bytes.Reader, d *decoder) (string, error) { return readBytes32Text(r) })
	addTextRecord(0xA4, "StartListText", func(r *bytes.Reader, d *decoder) (string, error) { return readListText(r, d) })
	addTextRecord(0xA6, "EndListText", func(r *bytes.Reader, d *decoder) (string, error) { return "", nil })
	addTextRecord(0xA8, "EmptyText", func(r *bytes.Reader, d *decoder) (string, error) { return "", nil })
	addTextRecord(0xAA, "DictionaryText", func(r *bytes.Reader, d *decoder) (string, error) { return readDictionaryString(r, d) })
	addTextRecord(0xAC, "UniqueIdText", func(r *bytes.Reader, d *decoder) (string, error) { return readUniqueIdText(r) })
	addTextRecord(0xAE, "TimeSpanText", func(r *bytes.Reader, d *decoder) (string, error) { return readTimeSpanText(r) })
	addTextRecord(0xB0, "UuidText", func(r *bytes.Reader, d *decoder) (string, error) { return readUuidText(r) })
	addTextRecord(0xB2, "UuidText", func(r *bytes.Reader, d *decoder) (string, error) { return readUInt64Text(r) })
	addTextRecord(0xB4, "BoolText", func(r *bytes.Reader, d *decoder) (string, error) { return readBoolText(r) })
	addTextRecord(0xB6, "UnicodeChars8Text", func(r *bytes.Reader, d *decoder) (string, error) { return readUnicodeChars8Text(r) })
	addTextRecord(0xB8, "UnicodeChars16Text", func(r *bytes.Reader, d *decoder) (string, error) { return readUnicodeChars16Text(r) })
	addTextRecord(0xBA, "UnicodeChars32Text", func(r *bytes.Reader, d *decoder) (string, error) { return readUnicodeChars32Text(r) })
	addTextRecord(0xBC, "QNameDictionaryText", func(r *bytes.Reader, d *decoder) (string, error) { return readQNameDictionaryText(r, d) })

	addAzRecords()
}

func addTextRecord(recordId byte, textName string, charData func(*bytes.Reader, *decoder) (string, error)) {
	records[recordId] = func(decoder *decoder) record {
		return &textRecordBase{decoder: decoder, withEndElement: false, textName: textName, recordId: recordId, charData: charData}
	}
	records[recordId+1] = func(decoder *decoder) record {
		return &textRecordBase{decoder: decoder, withEndElement: true, textName: textName + "WithEndElement", recordId: recordId + 1, charData: charData}
	}
}

func addAzRecords() {
	for b := 0; b < 26; b++ {
		byt := byte(b)
		records[byte(0x0C+byt)] = func(decoder *decoder) record { return &prefixDictionaryAttributeAZRecord{&attributeRecordBase{decoder: decoder}, byt, 0} }
		records[byte(0x26+byt)] = func(decoder *decoder) record { return &prefixAttributeAZRecord{&attributeRecordBase{decoder: decoder}, byt} }
		records[byte(0x44+byt)] = func(decoder *decoder) record { return &prefixDictionaryElementAZRecord{&elementRecordBase{decoder: decoder}, byt, 0} }
		records[byte(0x5E+byt)] = func(decoder *decoder) record { return &prefixElementAZRecord{&elementRecordBase{decoder: decoder}, "", byt} }
	}
}

func init() {
	initRecords()
}

//(0x01)
type endElementRecord struct {
	*elementRecordBase
}

func (r *endElementRecord) getName() string {
	return "EndElementRecord (0x01)"
}

func (r *endElementRecord) decodeElement(x *xml.Encoder, reader *bytes.Reader) (record, error) {
	item := r.decoder.elementStack.Pop()
	element := item.(xml.StartElement)
	endElementToken := xml.EndElement{Name:xml.Name{Local:element.Name.Local,Space:element.Name.Space}}
	err := x.EncodeToken(endElementToken)
	return nil, err
}

func (r *endElementRecord) write(writer io.Writer) error {
	return errors.New("NotImplemented: endElementRecord.write")
}

//(0x04)
type shortAttributeRecord struct {
	*attributeRecordBase
}

func (r *shortAttributeRecord) getName() string {
	return "ShortAttributeRecord (0x04)"
}

func (r *shortAttributeRecord) decodeAttribute(x *xml.Encoder, reader *bytes.Reader) (xml.Attr, error) {
	name, err := readString(reader)
	if err != nil {
		return xml.Attr{}, err
	}
	record, err := getNextRecord(r.decoder, reader)
	if err != nil {
		return xml.Attr{}, err
	}
	textReader := record.(textRecordDecoder)
	text, err := textReader.readText(reader)
	if err != nil {
		return xml.Attr{}, err
	}
	return xml.Attr{Name: xml.Name{Local: name}, Value: text}, nil
}

//(0x05)
type attributeRecord struct {
	*attributeRecordBase
}

func (r *attributeRecord) getName() string {
	return "AttributeRecord (0x05)"
}

func (r *attributeRecord) decodeAttribute(x *xml.Encoder, reader *bytes.Reader) (xml.Attr, error) {
	prefix, err := readString(reader)
	if err != nil {
		return xml.Attr{}, err
	}
	name, err := readString(reader)
	if err != nil {
		return xml.Attr{}, err
	}
	rec, err := getNextRecord(r.decoder, reader)
	if err != nil {
		return xml.Attr{}, err
	}
	textReader := rec.(textRecordDecoder)
	text, err := textReader.readText(reader)
	if err != nil {
		return xml.Attr{}, err
	}
	return xml.Attr{Name: xml.Name{Local: prefix + ":" + name}, Value: text}, nil
}

//(0x06)
type shortDictionaryAttributeRecord struct {
	*attributeRecordBase
	nameIndex uint32
}

func (r *shortDictionaryAttributeRecord) getName() string {
	return "ShortDictionaryAttributeRecord (0x06)"
}

func (r *shortDictionaryAttributeRecord) decodeAttribute(x *xml.Encoder, reader *bytes.Reader) (xml.Attr, error) {
	name, err := readDictionaryString(reader, r.decoder)
	if err != nil {
		return xml.Attr{}, err
	}
	rec, err := getNextRecord(r.decoder, reader)
	if err != nil {
		return xml.Attr{}, err
	}
	textReader := rec.(textRecordDecoder)
	text, err := textReader.readText(reader)
	if err != nil {
		return xml.Attr{}, err
	}
	return xml.Attr{Name: xml.Name{Local: name}, Value: text}, nil
}

//(0x07)
type dictionaryAttributeRecord struct {
	*attributeRecordBase
}

func (r *dictionaryAttributeRecord) getName() string {
	return "DictionaryAttributeRecord (0x07)"
}

func (r *dictionaryAttributeRecord) decodeAttribute(x *xml.Encoder, reader *bytes.Reader) (xml.Attr, error) {
	prefix, err := readString(reader)
	if err != nil {
		return xml.Attr{}, err
	}
	name, err := readDictionaryString(reader, r.decoder)
	if err != nil {
		return xml.Attr{}, err
	}
	rec, err := getNextRecord(r.decoder, reader)
	if err != nil {
		return xml.Attr{}, err
	}
	textReader := rec.(textRecordDecoder)
	text, err := textReader.readText(reader)
	if err != nil {
		return xml.Attr{}, err
	}
	return xml.Attr{Name: xml.Name{Local: prefix + ":" + name}, Value: text}, nil
}

//(0x08)
type shortXmlnsAttributeRecord struct {
	*attributeRecordBase
}

func (r *shortXmlnsAttributeRecord) getName() string {
	return "ShortXmlnsAttributeRecord (0x08)"
}

func (r *shortXmlnsAttributeRecord) decodeAttribute(x *xml.Encoder, reader *bytes.Reader) (xml.Attr, error) {
	name := "xmlns"
	val, err := readString(reader)
	if err != nil {
		return xml.Attr{}, err
	}
	return xml.Attr{Name: xml.Name{Local: name}, Value: val}, nil
}

//(0x09)
type xmlnsAttributeRecord struct {
	*attributeRecordBase
}

func (r *xmlnsAttributeRecord) getName() string {
	return "XmlnsAttributeRecord (0x09)"
}

func (r *xmlnsAttributeRecord) decodeAttribute(x *xml.Encoder, reader *bytes.Reader) (xml.Attr, error) {
	prefix := "xmlns"
	name, err := readString(reader)
	if err != nil {
		return xml.Attr{}, err
	}
	val, err := readString(reader)
	if err != nil {
		return xml.Attr{}, err
	}
	return xml.Attr{Name: xml.Name{Local: prefix + ":" + name}, Value: val}, nil
}

//(0x0A)
type shortDictionaryXmlnsAttributeRecord struct {
	*attributeRecordBase
}

func (r *shortDictionaryXmlnsAttributeRecord) getName() string {
	return "ShortXmlnsAttributeRecord (0x0A)"
}

func (r *shortDictionaryXmlnsAttributeRecord) decodeAttribute(x *xml.Encoder, reader *bytes.Reader) (xml.Attr, error) {
	name := "xmlns"
	val, err := readDictionaryString(reader, r.decoder)
	if err != nil {
		return xml.Attr{}, err
	}
	return xml.Attr{Name: xml.Name{Local: name}, Value: val}, nil
}

//(0x0C-0x25)
type prefixDictionaryAttributeAZRecord struct {
	*attributeRecordBase
	prefixIndex byte
	nameIndex   uint32
}

func (r *prefixDictionaryAttributeAZRecord) getName() string {
	return fmt.Sprintf("PrefixDictionaryAttributeAZRecord (%#x)", byte(0x0C+r.prefixIndex))
}

func (r *prefixDictionaryAttributeAZRecord) decodeAttribute(x *xml.Encoder, reader *bytes.Reader) (xml.Attr, error) {
	name, err := readDictionaryString(reader, r.decoder)
	if err != nil {
		return xml.Attr{}, err
	}
	record, err := getNextRecord(r.decoder, reader)
	if err != nil {
		return xml.Attr{}, err
	}
	textRecord := record.(textRecordDecoder)
	if textRecord == nil {
		return xml.Attr{}, errors.New("Expected TextRecord")
	}
	text, err := textRecord.readText(reader)
	if err != nil {
		return xml.Attr{}, err
	}
	return xml.Attr{Name: xml.Name{Local: string('a'+r.prefixIndex) + ":" + name}, Value: text}, nil
}

//(0x26-0x3F)
type prefixAttributeAZRecord struct {
	*attributeRecordBase
	prefixIndex byte
}

func (r *prefixAttributeAZRecord) getName() string {
	return fmt.Sprintf("PrefixAttributeAZRecord (%#x)", byte(0x26+r.prefixIndex))
}

func (r *prefixAttributeAZRecord) decodeAttribute(x *xml.Encoder, reader *bytes.Reader) (xml.Attr, error) {
	name, err := readString(reader)
	if err != nil {
		return xml.Attr{}, err
	}
	record, err := getNextRecord(r.decoder, reader)
	if err != nil {
		return xml.Attr{}, err
	}
	textRecord := record.(textRecordDecoder)
	if textRecord == nil {
		return xml.Attr{}, errors.New("Expected TextRecord")
	}
	text, err := textRecord.readText(reader)
	if err != nil {
		return xml.Attr{}, err
	}
	return xml.Attr{Name: xml.Name{Local: string('a'+r.prefixIndex) + ":" + name}, Value: text}, nil
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

func (r *prefixDictionaryElementAZRecord) decodeElement(x *xml.Encoder, reader *bytes.Reader) (record, error) {
	name, err := readDictionaryString(reader, r.decoder)
	if err != nil {
		return nil, err
	}
	element := xml.StartElement{Name: xml.Name{Local: string('a'+r.prefixIndex) + ":" + name}}

	return r.readElementAttributes(element, x, reader)
}

func (r *prefixDictionaryElementAZRecord) write(writer io.Writer) error {
	writer.Write([]byte{0x44 + r.prefixIndex})
	_, err := writeMultiByteInt31(writer, r.nameIndex)
	return err
}

//(0x42)
type shortDictionaryElementRecord struct {
	*elementRecordBase
}

func (r *shortDictionaryElementRecord) getName() string {
	return "ShortDictionaryElement (0x42)"
}

func (r *shortDictionaryElementRecord) decodeElement(x *xml.Encoder, reader *bytes.Reader) (record, error) {
	name, err := readDictionaryString(reader, r.decoder)
	if err != nil {
		return nil, err
	}
	element := xml.StartElement{Name: xml.Name{Local: name}}

	return r.readElementAttributes(element, x, reader)
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

func (r *dictionaryElementRecord) decodeElement(x *xml.Encoder, reader *bytes.Reader) (record, error) {
	prefix, err := readString(reader)
	if err != nil {
		return nil, err
	}
	name, err := readDictionaryString(reader, r.decoder)
	if err != nil {
		return nil, err
	}
	element := xml.StartElement{Name: xml.Name{Local: prefix + ":" + name}}

	return r.readElementAttributes(element, x, reader)
}

//(0x0B)
type dictionaryXmlnsAttributeRecord struct {
	*attributeRecordBase
}

func (r *dictionaryXmlnsAttributeRecord) getName() string {
	return "dictionaryXmlnsAttribute (0x0B)"
}

func (r *dictionaryXmlnsAttributeRecord) decodeAttribute(x *xml.Encoder, reader *bytes.Reader) (xml.Attr, error) {
	name, err := readString(reader)
	if err != nil {
		return xml.Attr{}, err
	}

	val, err := readDictionaryString(reader, r.decoder)
	if err != nil {
		return xml.Attr{}, err
	}

	return xml.Attr{Name: xml.Name{Local: "xmlns:" + name}, Value: val}, nil
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

func (r *shortElementRecord) decodeElement(x *xml.Encoder, reader *bytes.Reader) (record, error) {
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

func (r *prefixElementAZRecord) decodeElement(x *xml.Encoder, reader *bytes.Reader) (record, error) {
	name, err := readString(reader)
	if err != nil {
		return nil, err
	}
	element := xml.StartElement{Name: xml.Name{Local: string('a'+r.prefixIndex) + ":" + name}}

	return r.readElementAttributes(element, x, reader)
}

func (r *prefixElementAZRecord) writeElement(w io.Writer, x xml.Token) error {
	w.Write([]byte{0x5E + r.prefixIndex})
	e := x.(xml.StartElement)
	writeString(w, e.Name.Local)
	return nil
}

// 0x41
type elementRecord struct {
	*elementRecordBase
	name   string
	prefix string
}

func (r *elementRecord) getName() string {
	return "elementRecord (0x41)"
}

func (r *elementRecord) decodeElement(x *xml.Encoder, reader *bytes.Reader) (record, error) {
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

// 0x02
type commentRecord struct {
	*textRecordBase
	text   string
}

func (r *commentRecord) getName() string {
	return "commentRecord (0x02)"
}

func (r *commentRecord) decodeText(x *xml.Encoder, reader *bytes.Reader) (string, error) {
	text, err := readString(reader)
	if err != nil {
		return "", err
	}
	element := xml.Comment(text)

	err = x.EncodeToken(element)
	if err != nil {
		return "", err
	}
	return text, nil
}

//func (r *commentRecord) write(writer io.Writer) error {
//	return errors.New("NotImplemented: elementRecord.write")
//}

// 0x03
type arrayRecord struct {
	*elementRecordBase
}

func (r *arrayRecord) getName() string {
	return "arrayRecord (0x03)"
}

func (r *arrayRecord) decodeElement(x *xml.Encoder, reader *bytes.Reader) (record, error) {
	rec, err := getNextRecord(r.decoder, reader)
	if err != nil {
		return rec, err
	}
	if !rec.isElement() {
		return nil, errors.New("Element expected!")
	}
	elementDecoder := rec.(elementRecordDecoder)
	rec, err = elementDecoder.decodeElement(x, reader)
	if err != nil {
		return rec, err
	}
	valRec, err := getNextRecord(r.decoder, reader)
	if err != nil {
		return valRec, err
	}
	valDecoder := valRec.(textRecordDecoder)
	len, err := readMultiByteInt31(reader)
	if err != nil {
		return nil, err
	}
	var i uint32
	var startElement xml.StartElement
	for i = 0; i < len; i++ {
		//fmt.Println("LOOP", r.decoder.elementStack.top.value)
		if i == 0 {
			startElement = r.decoder.elementStack.top.value.(xml.StartElement)
		} else {
			err = x.EncodeToken(startElement)
			if err != nil {
				return nil, err
			}
			r.decoder.elementStack.Push(startElement)
		}
		//fmt.Println("DecodeText", r.decoder.elementStack.top.value)
		_, err = valDecoder.decodeText(x, reader)
		//fmt.Println("DecodeText2", r.decoder.elementStack.top.value)
		if err != nil {
			return nil, err
		}
		if i < len {
			r.decoder.elementStack.Push(startElement)
		}
	}
	return nil, nil
}

//func (r *arrayRecord) write(writer io.Writer) error {
//	return errors.New("NotImplemented: elementRecord.write")
//}
