package nbfx

import (
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

type elementRecordDecoder interface {
	decodeElement(d *decoder) (record, error)
}

type elementRecordWriter interface {
	writeElement(e *encoder) error
}

type attributeRecordDecoder interface {
	decodeAttribute(d *decoder) (xml.Attr, error)
}

type textRecordDecoder interface {
	decodeText(d *decoder) (string, error)
	readText(d *decoder) (string, error)
}

type textRecordWriter interface {
	writeText(e *encoder) error
}

type recordBase struct {
	name       string
	id       byte
}

func (r *recordBase) isElement() bool   { return false }
func (r *recordBase) isAttribute() bool { return false }
func (r recordBase) getName() string { return r.name }

type elementRecordBase struct {
	recordBase
}

func (e *elementRecordBase) isElement() bool { return true }

func (r *elementRecordBase) readElementAttributes(element xml.StartElement, d *decoder) (record, error) {
	var peekRecord record
	var attributeToken xml.Attr

	// get next record
	//fmt.Println("getting next record")
	rec, err := getNextRecord(d)
	for err == nil && rec != nil {
		//fmt.Println("Processing record", rec.getName())
		if err != nil {
			return nil, err
		}

		var attrReader attributeRecordDecoder
		if rec.isAttribute() {
			attrReader = rec.(attributeRecordDecoder)

			attributeToken, err = attrReader.decodeAttribute(d)
			if err != nil {
				return nil, err
			}
			element.Attr = append(element.Attr, attributeToken)

			rec, err = getNextRecord(d)
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

	err = d.xml.EncodeToken(element)
	if err != nil {
		return nil, err
	}

	d.elementStack.Push(element)

	return peekRecord, nil
}

type attributeRecordBase struct{
	recordBase
}

func (r *attributeRecordBase) isAttribute() bool { return true }

type textRecordBase struct {
	recordBase
	withEndElement bool
	charData       func(d *decoder) (string, error)
}

func (r *textRecordBase) isElement() bool   { return false }
func (r *textRecordBase) isAttribute() bool { return false }

func (r *textRecordBase) getName() string {
	return fmt.Sprintf("%s (%#x)", r.name, r.id)
}

func (r *textRecordBase) readText(d *decoder) (string, error) {
	text, err := r.charData(d)
	if err != nil {
		return "", err
	}
	return text, nil
}

func (r *textRecordBase) decodeText(d *decoder) (string, error) {
	text, err := r.readText(d)
	if err != nil {
		return "", err
	}
	charData := xml.CharData([]byte(text))
	d.xml.EncodeToken(charData)
	if r.withEndElement {
		rec, err := getRecord(EndElement)
		endElementReader := rec.(elementRecordDecoder)
		if err == nil {
			endElementReader.decodeElement(d)
		}
	}
	return text, nil
}

func (r *textRecordBase) writeText(w *io.Writer) error {
	return errors.New("NotImplement: writeText for " + r.getName())
}

func getNextRecord(d *decoder) (record, error) {
	b, err := d.bin.ReadByte()
	if err != nil {
		return nil, err
	}

	return getRecord(b)
}

func getRecord(b byte) (record, error) {
	if rec, ok := records[b]; ok {
		return rec, nil
	}

	return nil, errors.New(fmt.Sprintf("Unknown record %#x", b))
}

var records = map[byte]record{}

func initRecords() {
	// Miscellaneous Records
	addRecord(EndElement, "EndElement", &endElementRecord{&elementRecordBase{}})
	records[Comment] = &commentRecord{&textRecordBase{}}
	records[Array] = &arrayRecord{&elementRecordBase{}}

	// Attribute Records
	records[ShortAttribute] = &shortAttributeRecord{&attributeRecordBase{}}
	records[Attribute] = &attributeRecord{&attributeRecordBase{}}
	records[ShortDictionaryAttribute] = &shortDictionaryAttributeRecord{&attributeRecordBase{}}
	records[DictionaryAttribute] = &dictionaryAttributeRecord{&attributeRecordBase{}}
	records[ShortXmlnsAttribute] = &shortXmlnsAttributeRecord{&attributeRecordBase{}}
	records[XmlnsAttribute] = &xmlnsAttributeRecord{&attributeRecordBase{}}
	records[ShortDictionaryXmlnsAttribute] = &shortDictionaryXmlnsAttributeRecord{&attributeRecordBase{}}
	records[DictionaryXmlnsAttribute] = &dictionaryXmlnsAttributeRecord{&attributeRecordBase{}}
	// PrefixDictionaryAttributeAZRecord ADDED IN addAzRecords()
	// PrefixAttributeAZRecord ADDED IN addAzRecords()

	// Element Records
	records[ShortElement] = &shortElementRecord{&elementRecordBase{}}
	records[Element] = &elementRecord{&elementRecordBase{}}
	records[ShortDictionaryElement] = &shortDictionaryElementRecord{&elementRecordBase{}}
	records[DictionaryElement] = &dictionaryElementRecord{&elementRecordBase{}}
	// PrefixDictionaryElementAZRecord ADDED IN addAzRecords()
	// PrefixElementAZRecord ADDED IN addAzRecords()

	// Text Records
	addTextRecord(ZeroText, "ZeroText", func(d *decoder) (string, error) { return "0", nil })
	addTextRecord(OneText, "OneText", func(d *decoder) (string, error) { return "1", nil })
	addTextRecord(FalseText, "FalseText", func(d *decoder) (string, error) { return "false", nil })
	addTextRecord(TrueText, "TrueText", func(d *decoder) (string, error) { return "true", nil})
	addTextRecord(Int8Text, "Int8Text", func(d *decoder) (string, error) { return readInt8Text(d) })
	addTextRecord(Int16Text, "Int16Text", func(d *decoder) (string, error) { return readInt16Text(d) })
	addTextRecord(Int32Text, "Int32Text", func(d *decoder) (string, error) { return readInt32Text(d) })
	addTextRecord(Int64Text, "Int64Text", func(d *decoder) (string, error) { return readInt64Text(d) })
	addTextRecord(FloatText, "FloatText", func(d *decoder) (string, error) { return readFloatText(d) })
	addTextRecord(DoubleText, "DoubleText", func(d *decoder) (string, error) { return readDoubleText(d) })
	addTextRecord(DecimalText, "DecimalText", func(d *decoder) (string, error) { return readDecimalText(d) })
	addTextRecord(DateTimeText, "DateTimeText", func(d *decoder) (string, error) { return readDateTimeText(d) })
	addTextRecord(Chars8Text, "Chars8Text", func(d *decoder) (string, error) { return readChars8Text(d) })
	addTextRecord(Chars16Text, "Chars16Text", func(d *decoder) (string, error) { return readChars16Text(d) })
	addTextRecord(Chars32Text, "Chars32Text", func(d *decoder) (string, error) { return readChars32Text(d) })
	addTextRecord(Bytes8Text, "Bytes8Text", func(d *decoder) (string, error) { return readBytes8Text(d) })
	addTextRecord(Bytes16Text, "Bytes16Text", func(d *decoder) (string, error) { return readBytes16Text(d) })
	addTextRecord(Bytes32Text, "Bytes32Text", func(d *decoder) (string, error) { return readBytes32Text(d) })
	addTextRecord(StartListText, "StartListText", func(d *decoder) (string, error) { return readListText(d) })
	addTextRecord(EndListText, "EndListText", func(d *decoder) (string, error) { return "", nil })
	addTextRecord(EmptyText, "EmptyText", func(d *decoder) (string, error) { return "", nil })
	addTextRecord(DictionaryText, "DictionaryText", func(d *decoder) (string, error) { return readDictionaryString(d) })
	addTextRecord(UniqueIdText, "UniqueIdText", func(d *decoder) (string, error) { return readUniqueIdText(d) })
	addTextRecord(TimeSpanText, "TimeSpanText", func(d *decoder) (string, error) { return readTimeSpanText(d) })
	addTextRecord(UuidText, "UuidText", func(d *decoder) (string, error) { return readUuidText(d) })
	addTextRecord(UInt64Text, "UInt64Text", func(d *decoder) (string, error) { return readUInt64Text(d) })
	addTextRecord(BoolText, "BoolText", func(d *decoder) (string, error) { return readBoolText(d) })
	addTextRecord(UnicodeChars8Text, "UnicodeChars8Text", func(d *decoder) (string, error) { return readUnicodeChars8Text(d) })
	addTextRecord(UnicodeChars16Text, "UnicodeChars16Text", func(d *decoder) (string, error) { return readUnicodeChars16Text(d) })
	addTextRecord(UnicodeChars32Text, "UnicodeChars32Text", func(d *decoder) (string, error) { return readUnicodeChars32Text(d) })
	addTextRecord(QNameDictionaryText, "QNameDictionaryText", func(d *decoder) (string, error) { return readQNameDictionaryText(d) })

	addAzRecords()
}

func addRecord(id byte, name string, rec record) {
	base := rec.(recordBase)
	base.id = id
	base.name = name
	records[id] = rec
}

func addTextRecord(recordId byte, textName string, charData func(*decoder) (string, error)) {
	records[recordId] = &textRecordBase{recordBase{name: textName, id: recordId}, false, charData}
	records[recordId+1] = &textRecordBase{recordBase{name: textName + "WithEndElement", id: recordId + 1}, true, charData}
}

func addAzRecords() {
	for b := 0; b < 26; b++ {
		byt := byte(b)
		records[byte(PrefixDictionaryAttributeA+byt)] = &prefixDictionaryAttributeAZRecord{&attributeRecordBase{}, byt, 0}
		records[byte(PrefixAttributeA+byt)] = &prefixAttributeAZRecord{&attributeRecordBase{}, byt}
		records[byte(PrefixDictionaryElementA+byt)] = &prefixDictionaryElementAZRecord{&elementRecordBase{}, byt, 0}
		records[byte(PrefixElementA+byt)] = &prefixElementAZRecord{&elementRecordBase{}, "", byt}
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

func (r *endElementRecord) decodeElement(d *decoder) (record, error) {
	item := d.elementStack.Pop()
	element := item.(xml.StartElement)
	endElementToken := xml.EndElement{Name:xml.Name{Local:element.Name.Local,Space:element.Name.Space}}
	err := d.xml.EncodeToken(endElementToken)
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

func (r *shortAttributeRecord) decodeAttribute(d *decoder) (xml.Attr, error) {
	name, err := readString(d.bin)
	if err != nil {
		return xml.Attr{}, err
	}
	record, err := getNextRecord(d)
	if err != nil {
		return xml.Attr{}, err
	}
	textReader := record.(textRecordDecoder)
	text, err := textReader.readText(d)
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

func (r *attributeRecord) decodeAttribute(d *decoder) (xml.Attr, error) {
	prefix, err := readString(d.bin)
	if err != nil {
		return xml.Attr{}, err
	}
	name, err := readString(d.bin)
	if err != nil {
		return xml.Attr{}, err
	}
	rec, err := getNextRecord(d)
	if err != nil {
		return xml.Attr{}, err
	}
	textReader := rec.(textRecordDecoder)
	text, err := textReader.readText(d)
	if err != nil {
		return xml.Attr{}, err
	}
	return xml.Attr{Name: xml.Name{Local: prefix + ":" + name}, Value: text}, nil
}

//(0x06)
type shortDictionaryAttributeRecord struct {
	*attributeRecordBase
}

func (r *shortDictionaryAttributeRecord) getName() string {
	return "ShortDictionaryAttributeRecord (0x06)"
}

func (r *shortDictionaryAttributeRecord) decodeAttribute(d *decoder) (xml.Attr, error) {
	name, err := readDictionaryString(d)
	if err != nil {
		return xml.Attr{}, err
	}
	rec, err := getNextRecord(d)
	if err != nil {
		return xml.Attr{}, err
	}
	textReader := rec.(textRecordDecoder)
	text, err := textReader.readText(d)
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

func (r *dictionaryAttributeRecord) decodeAttribute(d *decoder) (xml.Attr, error) {
	prefix, err := readString(d.bin)
	if err != nil {
		return xml.Attr{}, err
	}
	name, err := readDictionaryString(d)
	if err != nil {
		return xml.Attr{}, err
	}
	rec, err := getNextRecord(d)
	if err != nil {
		return xml.Attr{}, err
	}
	textReader := rec.(textRecordDecoder)
	text, err := textReader.readText(d)
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

func (r *shortXmlnsAttributeRecord) decodeAttribute(d *decoder) (xml.Attr, error) {
	name := "xmlns"
	val, err := readString(d.bin)
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

func (r *xmlnsAttributeRecord) decodeAttribute(d *decoder) (xml.Attr, error) {
	prefix := "xmlns"
	name, err := readString(d.bin)
	if err != nil {
		return xml.Attr{}, err
	}
	val, err := readString(d.bin)
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

func (r *shortDictionaryXmlnsAttributeRecord) decodeAttribute(d *decoder) (xml.Attr, error) {
	name := "xmlns"
	val, err := readDictionaryString(d)
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

func (r *prefixDictionaryAttributeAZRecord) decodeAttribute(d *decoder) (xml.Attr, error) {
	name, err := readDictionaryString(d)
	if err != nil {
		return xml.Attr{}, err
	}
	record, err := getNextRecord(d)
	if err != nil {
		return xml.Attr{}, err
	}
	textRecord := record.(textRecordDecoder)
	if textRecord == nil {
		return xml.Attr{}, errors.New("Expected TextRecord")
	}
	text, err := textRecord.readText(d)
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

func (r *prefixAttributeAZRecord) decodeAttribute(d *decoder) (xml.Attr, error) {
	name, err := readString(d.bin)
	if err != nil {
		return xml.Attr{}, err
	}
	record, err := getNextRecord(d)
	if err != nil {
		return xml.Attr{}, err
	}
	textRecord := record.(textRecordDecoder)
	if textRecord == nil {
		return xml.Attr{}, errors.New("Expected TextRecord")
	}
	text, err := textRecord.readText(d)
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

func (r *prefixDictionaryElementAZRecord) decodeElement(d *decoder) (record, error) {
	name, err := readDictionaryString(d)
	if err != nil {
		return nil, err
	}
	element := xml.StartElement{Name: xml.Name{Local: string('a'+r.prefixIndex) + ":" + name}}

	return r.readElementAttributes(element, d)
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

func (r *shortDictionaryElementRecord) decodeElement(d *decoder) (record, error) {
	name, err := readDictionaryString(d)
	if err != nil {
		return nil, err
	}
	element := xml.StartElement{Name: xml.Name{Local: name}}

	return r.readElementAttributes(element, d)
}

//(0x43)
type dictionaryElementRecord struct {
	*elementRecordBase
}

func (r *dictionaryElementRecord) getName() string {
	return "DictionaryElement (0x43)"
}

func (r *dictionaryElementRecord) decodeElement(d *decoder) (record, error) {
	prefix, err := readString(d.bin)
	if err != nil {
		return nil, err
	}
	name, err := readDictionaryString(d)
	if err != nil {
		return nil, err
	}
	element := xml.StartElement{Name: xml.Name{Local: prefix + ":" + name}}

	return r.readElementAttributes(element, d)
}

//(0x0B)
type dictionaryXmlnsAttributeRecord struct {
	*attributeRecordBase
}

func (r *dictionaryXmlnsAttributeRecord) getName() string {
	return "dictionaryXmlnsAttribute (0x0B)"
}

func (r *dictionaryXmlnsAttributeRecord) decodeAttribute(d *decoder) (xml.Attr, error) {
	name, err := readString(d.bin)
	if err != nil {
		return xml.Attr{}, err
	}

	val, err := readDictionaryString(d)
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
}

func (r *shortElementRecord) getName() string {
	return "shortElementRecord (0x40)"
}

func (r *shortElementRecord) decodeElement(d *decoder) (record, error) {
	name, err := readString(d.bin)
	if err != nil {
		return nil, err
	}
	element := xml.StartElement{Name: xml.Name{Local: name}}

	return r.readElementAttributes(element, d)
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

func (r *prefixElementAZRecord) decodeElement(d *decoder) (record, error) {
	name, err := readString(d.bin)
	if err != nil {
		return nil, err
	}
	element := xml.StartElement{Name: xml.Name{Local: string('a'+r.prefixIndex) + ":" + name}}

	return r.readElementAttributes(element, d)
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
}

func (r *elementRecord) getName() string {
	return "elementRecord (0x41)"
}

func (r *elementRecord) decodeElement(d *decoder) (record, error) {
	prefix, err := readString(d.bin)
	if err != nil {
		return nil, err
	}
	name, err := readString(d.bin)
	if err != nil {
		return nil, err
	}
	element := xml.StartElement{Name: xml.Name{Local: prefix + ":" + name}}

	return r.readElementAttributes(element, d)
}

//func (r *elementRecord) write(writer io.Writer) error {
//	return errors.New("NotImplemented: elementRecord.write")
//}

// 0x02
type commentRecord struct {
	*textRecordBase
}

func (r *commentRecord) getName() string {
	return "commentRecord (0x02)"
}

func (r *commentRecord) decodeText(d *decoder) (string, error) {
	text, err := readString(d.bin)
	if err != nil {
		return "", err
	}
	element := xml.Comment(text)

	err = d.xml.EncodeToken(element)
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

func (r *arrayRecord) decodeElement(d *decoder) (record, error) {
	rec, err := getNextRecord(d)
	if err != nil {
		return rec, err
	}
	if !rec.isElement() {
		return nil, errors.New("Element expected!")
	}
	elementDecoder := rec.(elementRecordDecoder)
	rec, err = elementDecoder.decodeElement(d)
	if err != nil {
		return rec, err
	}
	valRec, err := getNextRecord(d)
	if err != nil {
		return valRec, err
	}
	valDecoder := valRec.(textRecordDecoder)
	len, err := readMultiByteInt31(d.bin)
	if err != nil {
		return nil, err
	}
	var i uint32
	var startElement xml.StartElement
	for i = 0; i < len; i++ {
		//fmt.Println("LOOP", r.decoder.elementStack.top.value)
		if i == 0 {
			startElement = d.elementStack.top.value.(xml.StartElement)
		} else {
			err = d.xml.EncodeToken(startElement)
			if err != nil {
				return nil, err
			}
			d.elementStack.Push(startElement)
		}
		//fmt.Println("DecodeText", r.decoder.elementStack.top.value)
		_, err = valDecoder.decodeText(d)
		//fmt.Println("DecodeText2", r.decoder.elementStack.top.value)
		if err != nil {
			return nil, err
		}
		if i < len {
			d.elementStack.Push(startElement)
		}
	}
	return nil, nil
}

//func (r *arrayRecord) write(writer io.Writer) error {
//	return errors.New("NotImplemented: elementRecord.write")
//}
