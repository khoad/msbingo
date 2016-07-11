package nbfx

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

type record interface {
	isStartElement() bool
	isEndElement() bool
	isAttribute() bool
	isText() bool
	getName() string
}

type elementRecordDecoder interface {
	decodeElement(d *decoder) (record, error)
}

type elementRecordEncoder interface {
	encodeElement(e *encoder, element xml.StartElement) error
}

type attributeRecordDecoder interface {
	decodeAttribute(d *decoder) (xml.Attr, error)
}

type attributeRecordEncoder interface {
	encodeAttribute(e *encoder, attr xml.Attr) error
}

type textRecordDecoder interface {
	decodeText(d *decoder, trd textRecordDecoder) (string, error)
	readText(d *decoder) (string, error)
}

type textRecordEncoder interface {
	encodeText(e *encoder, tre textRecordEncoder, text string) error
	writeText(e *encoder, text string) error
}

type recordBase struct {
	name       string
	id       byte
}

func (r *recordBase) isStartElement() bool   { return false }
func (r *recordBase) isEndElement() bool { return false }
func (r *recordBase) isAttribute() bool { return false }
func (r *recordBase) isText() bool { return false }
func (r *recordBase) getName() string { return r.name }

type elementRecordBase struct {
	recordBase
}

func (e *elementRecordBase) isStartElement() bool { return true }

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

	err = d.xml.EncodeToken(element)
	if err != nil {
		return nil, err
	}

	d.elementStack.Push(element)

	return peekRecord, nil
}

func (r *elementRecordBase) encodeAttributes(e *encoder, attrs []xml.Attr) error {
	var err error
	for _, attr := range attrs {
		err = r.encodeAttribute(e, attr)
		if err != nil {
			break
		}
	}
	return err
}

func (r *elementRecordBase) encodeAttribute(e *encoder, attr xml.Attr) error {
	rec, err := e.getAttributeRecordFromToken(attr)
	if err != nil {
		return err
	}
	attrRec := rec.(attributeRecordEncoder)
	return attrRec.encodeAttribute(e, attr)
}

func (r *elementRecordBase) encodeElement(e *encoder, element xml.StartElement) error {
	return errors.New(fmt.Sprint("NotImplemented: encodeElement on", r))
}

type attributeRecordBase struct{
	recordBase
}

func (r *attributeRecordBase) isAttribute() bool { return true }

func (r *attributeRecordBase) encodeAttribute(e *encoder, attr xml.Attr) error {
	return errors.New(fmt.Sprint("NotImplemented: encodeAttribute on", r))
}

type textRecordBase struct {
	recordBase
	withEndElement bool
}

func (r *textRecordBase) isText() bool {
	return true
}

func (r *textRecordBase) readText(d *decoder) (string, error) {
	return "", errors.New(fmt.Sprintf("NotImplemented: %v.readText %v", r.name, r))
}

func (r *textRecordBase) writeText(e *encoder, text string) error {
	return errors.New("NotImplemented: " + r.name + ".writeText")
}

func (r *textRecordBase) decodeText(d *decoder, trd textRecordDecoder) (string, error) {
	// This is an ugly hack to allow polymorphic readText behavior
	text, err := trd.readText(d)
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

func (r *textRecordBase) encodeText(e *encoder, tre textRecordEncoder, text string) error {
	err := e.bin.WriteByte(r.id)
	if err != nil {
		return err
	}

	err = tre.writeText(e, text)
	if err != nil {
		return err
	}

	return nil
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
	records[EndElement] = &endElementRecord{elementRecordBase{recordBase{"EndElement",EndElement}}}
	records[Comment] = &commentRecord{textRecordBase{recordBase{"Comment",Comment},false}}
	records[Array] = &arrayRecord{elementRecordBase{recordBase{"Array",Array}}}

	// Attribute Records
	records[ShortAttribute] = &shortAttributeRecord{attributeRecordBase{recordBase{"ShortAttribute",ShortAttribute}}}
	records[Attribute] = &attributeRecord{attributeRecordBase{recordBase{"Attribute",Attribute}}}
	records[ShortDictionaryAttribute] = &shortDictionaryAttributeRecord{attributeRecordBase{recordBase{"ShortDictionaryAttribute",ShortDictionaryAttribute}}}
	records[DictionaryAttribute] = &dictionaryAttributeRecord{attributeRecordBase{recordBase{"DictionaryAttribute",DictionaryAttribute}}}
	records[ShortXmlnsAttribute] = &shortXmlnsAttributeRecord{attributeRecordBase{recordBase{"ShortXmlnsAttribute",ShortXmlnsAttribute}}}
	records[XmlnsAttribute] = &xmlnsAttributeRecord{attributeRecordBase{recordBase{"XmlnsAttribute",XmlnsAttribute}}}
	records[ShortDictionaryXmlnsAttribute] = &shortDictionaryXmlnsAttributeRecord{attributeRecordBase{recordBase{"ShortDictionaryXmlnsAttribute",ShortDictionaryXmlnsAttribute}}}
	records[DictionaryXmlnsAttribute] = &dictionaryXmlnsAttributeRecord{attributeRecordBase{recordBase{"DictionaryXmlnsAttribute",DictionaryXmlnsAttribute}}}
	// PrefixDictionaryAttributeAZRecord ADDED IN addAzRecords()
	// PrefixAttributeAZRecord ADDED IN addAzRecords()

	// Element Records
	records[ShortElement] = &shortElementRecord{elementRecordBase{recordBase{"ShortElement",ShortElement}}}
	records[Element] = &elementRecord{elementRecordBase{recordBase{"Element",Element}}}
	records[ShortDictionaryElement] = &shortDictionaryElementRecord{elementRecordBase{recordBase{"ShortDictionaryElement",ShortDictionaryElement}}}
	records[DictionaryElement] = &dictionaryElementRecord{elementRecordBase{recordBase{"DictionaryElement",DictionaryElement}}}
	// PrefixDictionaryElementAZRecord ADDED IN addAzRecords()
	// PrefixElementAZRecord ADDED IN addAzRecords()

	// Text Records
	records[ZeroText] = &zeroTextRecord{textRecordBase{recordBase{"ZeroText", ZeroText},false}}
	records[ZeroTextWithEndElement] = &zeroTextRecord{textRecordBase{recordBase{"ZeroTextWithEndElement", ZeroTextWithEndElement},true}}
	records[OneText] = &oneTextRecord{textRecordBase{recordBase{"OneText", OneText},false}}
	records[OneTextWithEndElement] = &oneTextRecord{textRecordBase{recordBase{"OneTextWithEndElement", OneTextWithEndElement},true}}
	records[FalseText] = &falseTextRecord{textRecordBase{recordBase{"FalseText", FalseText},false}}
	records[FalseTextWithEndElement] = &falseTextRecord{textRecordBase{recordBase{"FalseTextWithEndElement", FalseTextWithEndElement},true}}
	records[TrueText] = &trueTextRecord{textRecordBase{recordBase{"TrueText", TrueText},false}}
	records[TrueTextWithEndElement] = &trueTextRecord{textRecordBase{recordBase{"TrueTextWithEndElement", TrueTextWithEndElement},true}}
	records[Int8Text] = &int8TextRecord{textRecordBase{recordBase{"Int8Text", Int8Text},false}}
	records[Int8TextWithEndElement] = &int8TextRecord{textRecordBase{recordBase{"Int8TextWithEndElement", Int8TextWithEndElement},true}}
	records[Int16Text] = &int16TextRecord{textRecordBase{recordBase{"Int16Text", Int16Text},false}}
	records[Int16TextWithEndElement] = &int16TextRecord{textRecordBase{recordBase{"Int16TextWithEndElement", Int16TextWithEndElement},true}}
	records[Int32Text] = &int32TextRecord{textRecordBase{recordBase{"Int32Text", Int32Text},false}}
	records[Int32TextWithEndElement] = &int32TextRecord{textRecordBase{recordBase{"Int32TextWithEndElement", Int32TextWithEndElement},true}}
	records[Int64Text] = &int64TextRecord{textRecordBase{recordBase{"Int64Text", Int64Text},false}}
	records[Int64TextWithEndElement] = &int64TextRecord{textRecordBase{recordBase{"Int64TextWithEndElement", Int64TextWithEndElement},true}}
	records[FloatText] = &floatTextRecord{textRecordBase{recordBase{"FloatText", FloatText},false}}
	records[FloatTextWithEndElement] = &floatTextRecord{textRecordBase{recordBase{"FloatTextWithEndElement", FloatTextWithEndElement},true}}
	records[DoubleText] = &doubleTextRecord{textRecordBase{recordBase{"DoubleText", DoubleText},false}}
	records[DoubleTextWithEndElement] = &doubleTextRecord{textRecordBase{recordBase{"DoubleTextWithEndElement", DoubleTextWithEndElement},true}}
	records[DecimalText] = &decimalTextRecord{textRecordBase{recordBase{"DecimalText", DecimalText},false}}
	records[DecimalTextWithEndElement] = &decimalTextRecord{textRecordBase{recordBase{"DecimalTextWithEndElement", DecimalTextWithEndElement},true}}
	records[DateTimeText] = &dateTimeTextRecord{textRecordBase{recordBase{"DateTimeText", DateTimeText},false}}
	records[DateTimeTextWithEndElement] = &dateTimeTextRecord{textRecordBase{recordBase{"DateTimeTextWithEndElement", DateTimeTextWithEndElement},true}}
	records[Chars8Text] = &chars8TextRecord{textRecordBase{recordBase{"Chars8Text", Chars8Text},false}}
	records[Chars8TextWithEndElement] = &chars8TextRecord{textRecordBase{recordBase{"Chars8TextWithEndElement", Chars8TextWithEndElement},true}}
	records[Chars16Text] = &chars16TextRecord{textRecordBase{recordBase{"Chars16Text", Chars16Text},false}}
	records[Chars16TextWithEndElement] = &chars16TextRecord{textRecordBase{recordBase{"Chars16TextWithEndElement", Chars16TextWithEndElement},true}}
	records[Chars32Text] = &chars32TextRecord{textRecordBase{recordBase{"Chars32Text", Chars32Text},false}}
	records[Chars32TextWithEndElement] = &chars32TextRecord{textRecordBase{recordBase{"Chars32TextWithEndElement", Chars32TextWithEndElement},true}}
	records[Bytes8Text] = &bytes8TextRecord{textRecordBase{recordBase{"Bytes8Text", Bytes8Text},false}}
	records[Bytes8TextWithEndElement] = &bytes8TextRecord{textRecordBase{recordBase{"Bytes8TextWithEndElement", Bytes8TextWithEndElement},true}}
	records[Bytes16Text] = &bytes16TextRecord{textRecordBase{recordBase{"Bytes16Text", Bytes16Text},false}}
	records[Bytes16TextWithEndElement] = &bytes16TextRecord{textRecordBase{recordBase{"Bytes16TextWithEndElement", Bytes16TextWithEndElement},true}}
	records[Bytes32Text] = &bytes32TextRecord{textRecordBase{recordBase{"Bytes32Text", Bytes32Text},false}}
	records[Bytes32TextWithEndElement] = &bytes32TextRecord{textRecordBase{recordBase{"Bytes32TextWithEndElement", Bytes32TextWithEndElement},true}}
	records[StartListText] = &startListTextRecord{textRecordBase{recordBase{"StartListText", StartListText},false}}
	records[StartListTextWithEndElement] = &startListTextRecord{textRecordBase{recordBase{"StartListTextWithEndElement", StartListTextWithEndElement},true}}
	records[EndListText] = &endListTextRecord{textRecordBase{recordBase{"EndListText", EndListText},false}}
	records[EndListTextWithEndElement] = &endListTextRecord{textRecordBase{recordBase{"EndListTextWithEndElement", EndListTextWithEndElement},true}}
	records[EmptyText] = &emptyTextRecord{textRecordBase{recordBase{"EmptyText", EmptyText},false}}
	records[EmptyTextWithEndElement] = &emptyTextRecord{textRecordBase{recordBase{"EmptyTextWithEndElement", EmptyTextWithEndElement},true}}
	records[DictionaryText] = &dictionaryTextRecord{textRecordBase{recordBase{"DictionaryText", DictionaryText},false}}
	records[DictionaryTextWithEndElement] = &dictionaryTextRecord{textRecordBase{recordBase{"DictionaryTextWithEndElement", DictionaryTextWithEndElement},true}}
	records[UniqueIdText] = &uniqueIdTextRecord{textRecordBase{recordBase{"UniqueIdText", UniqueIdText},false}}
	records[UniqueIdTextWithEndElement] = &uniqueIdTextRecord{textRecordBase{recordBase{"UniqueIdTextWithEndElement", UniqueIdTextWithEndElement},true}}
	records[TimeSpanText] = &timeSpanTextRecord{textRecordBase{recordBase{"TimeSpanText", TimeSpanText},false}}
	records[TimeSpanTextWithEndElement] = &timeSpanTextRecord{textRecordBase{recordBase{"TimeSpanTextWithEndElement", TimeSpanTextWithEndElement},true}}
	records[UuidText] = &uuidTextRecord{textRecordBase{recordBase{"UuidText", UuidText},false}}
	records[UuidTextWithEndElement] = &uuidTextRecord{textRecordBase{recordBase{"UuidTextWithEndElement", UuidTextWithEndElement},true}}
	records[UInt64Text] = &uInt64TextRecord{textRecordBase{recordBase{"UInt64Text", UInt64Text},false}}
	records[UInt64TextWithEndElement] = &uInt64TextRecord{textRecordBase{recordBase{"UInt64TextWithEndElement", UInt64TextWithEndElement},true}}
	records[BoolText] = &boolTextRecord{textRecordBase{recordBase{"BoolText", BoolText},false}}
	records[BoolTextWithEndElement] = &boolTextRecord{textRecordBase{recordBase{"BoolTextWithEndElement", BoolTextWithEndElement},true}}
	records[UnicodeChars8Text] = &unicodeChars8TextRecord{textRecordBase{recordBase{"UnicodeChars8Text", UnicodeChars8Text},false}}
	records[UnicodeChars8TextWithEndElement] = &unicodeChars8TextRecord{textRecordBase{recordBase{"UnicodeChars8TextWithEndElement", UnicodeChars8TextWithEndElement},true}}
	records[UnicodeChars16Text] = &unicodeChars16TextRecord{textRecordBase{recordBase{"UnicodeChars16Text", UnicodeChars16Text},false}}
	records[UnicodeChars16TextWithEndElement] = &unicodeChars16TextRecord{textRecordBase{recordBase{"UnicodeChars16TextWithEndElement", UnicodeChars16TextWithEndElement},true}}
	records[UnicodeChars32Text] = &unicodeChars32TextRecord{textRecordBase{recordBase{"UnicodeChars32Text", UnicodeChars32Text},false}}
	records[UnicodeChars32TextWithEndElement] = &unicodeChars32TextRecord{textRecordBase{recordBase{"UnicodeChars32TextWithEndElement", UnicodeChars32TextWithEndElement},true}}
	records[QNameDictionaryText] = &qNameDictionaryTextRecord{textRecordBase{recordBase{"QNameDictionaryText", QNameDictionaryText},false}}
	records[QNameDictionaryTextWithEndElement] = &qNameDictionaryTextRecord{textRecordBase{recordBase{"QNameDictionaryTextWithEndElement", QNameDictionaryTextWithEndElement},true}}

	addAzRecords()
}

func addAzRecords() {
	var b byte
	for b = 0; b < 26; b++ {
		byt := byte(b)
		prefix := strings.ToUpper(string('a'+b))
		records[byte(PrefixDictionaryAttributeA+byt)] = &prefixDictionaryAttributeAZRecord{attributeRecordBase{recordBase{"PrefixDictionaryAttribute"+prefix,byte(PrefixDictionaryAttributeA+b)}}}
		records[byte(PrefixAttributeA+byt)] = &prefixAttributeAZRecord{attributeRecordBase{recordBase{"PrefixAttribute"+prefix,byte(PrefixAttributeA+b)}}}
		records[byte(PrefixDictionaryElementA+byt)] = &prefixDictionaryElementAZRecord{elementRecordBase{recordBase{"PrefixDictionaryElement"+prefix,byte(PrefixDictionaryElementA+b)}}}
		records[byte(PrefixElementA+byt)] = &prefixElementAZRecord{elementRecordBase{recordBase{"PrefixElement"+prefix, byte(PrefixElementA+byt)}}}
	}
}

func init() {
	initRecords()
}

//(0x01)
type endElementRecord struct {
	elementRecordBase
}

func (r *endElementRecord) getName() string {
	return "EndElementRecord (0x01)"
}

func (r *endElementRecord) isStartElement() bool {
	return false
}

func (r *endElementRecord) isEndElement() bool {
	return true
}

func (r *endElementRecord) decodeElement(d *decoder) (record, error) {
	item := d.elementStack.Pop()
	element := item.(xml.StartElement)
	endElementToken := xml.EndElement{Name:xml.Name{Local:element.Name.Local,Space:element.Name.Space}}
	err := d.xml.EncodeToken(endElementToken)
	return nil, err
}

func (r *endElementRecord) encodeElement(e *encoder, element xml.StartElement) error {
	_, err := e.bin.Write([]byte{r.id})
	return err
}

//(0x04)
type shortAttributeRecord struct {
	attributeRecordBase
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
	attributeRecordBase
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
	attributeRecordBase
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
	attributeRecordBase
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
	attributeRecordBase
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
	attributeRecordBase
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
	attributeRecordBase
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
	attributeRecordBase
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
	return xml.Attr{Name: xml.Name{Local: string('a'+r.id-PrefixDictionaryAttributeA) + ":" + name}, Value: text}, nil
}

func (r *prefixDictionaryAttributeAZRecord) encodeAttribute(e *encoder, attr xml.Attr) error {
	err := e.bin.WriteByte(r.id)
	if err != nil {
		return err
	}
	_, err = writeDictionaryString(e, attr.Name.Local)
	if err != nil {
		return err
	}
	textRecord, err := e.getTextRecordFromText(attr.Value, false)
	//fmt.Println("prefixDictionaryAttributeAZRecord gotTextRecord", textRecord)
	if err != nil {
		return err
	}
	textEncoder := textRecord.(textRecordEncoder)
	return textEncoder.encodeText(e, textEncoder, attr.Value)
}

//(0x26-0x3F)
type prefixAttributeAZRecord struct {
	attributeRecordBase
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
	return xml.Attr{Name: xml.Name{Local: string('a'+r.id-PrefixAttributeA) + ":" + name}, Value: text}, nil
}

//(0x44-0x5D)
type prefixDictionaryElementAZRecord struct {
	elementRecordBase
}

func (r *prefixDictionaryElementAZRecord) decodeElement(d *decoder) (record, error) {
	name, err := readDictionaryString(d)
	if err != nil {
		return nil, err
	}
	element := xml.StartElement{Name: xml.Name{Local: string('a'+r.id-PrefixDictionaryElementA) + ":" + name}}

	return r.readElementAttributes(element, d)
}

func (r *prefixDictionaryElementAZRecord) encodeElement(e *encoder, element xml.StartElement) error {
	//fmt.Println("--->", element, e.bin)
	e.bin.Write([]byte{r.id})
	_, err := writeMultiByteInt31(e, e.dict[element.Name.Local])
	err = r.encodeAttributes(e, element.Attr)
	return err
}

//(0x42)
type shortDictionaryElementRecord struct {
	elementRecordBase
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
	elementRecordBase
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
	attributeRecordBase
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

func (r *dictionaryXmlnsAttributeRecord) encodeAttribute(e *encoder, attr xml.Attr) error {
	err := e.bin.WriteByte(r.id)
	if err != nil {
		return err
	}
	_, err = writeString(e, attr.Name.Local)
	if err != nil {
		return err
	}
	_, err = writeDictionaryString(e, attr.Value)
	return err
}

// 0x40
type shortElementRecord struct {
	elementRecordBase
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

func (r *shortElementRecord) encodeElement(e *encoder, element xml.StartElement) error {
	_, err := e.bin.Write([]byte{ShortElement})
	if err != nil {
		return err
	}
	_, err = writeString(e, element.Name.Local)
	return err
}

// 0x5E-0x77
type prefixElementAZRecord struct {
	elementRecordBase
}

func (r *prefixElementAZRecord) getName() string {
	return fmt.Sprint(r)
}

func (r *prefixElementAZRecord) decodeElement(d *decoder) (record, error) {
	name, err := readString(d.bin)
	if err != nil {
		return nil, err
	}
	element := xml.StartElement{Name: xml.Name{Local: string('a'+(r.id - PrefixElementA)) + ":" + name}}

	return r.readElementAttributes(element, d)
}

func (r *prefixElementAZRecord) encodeElement(e *encoder, element xml.StartElement) error {
	_, err := e.bin.Write([]byte{PrefixElementA + (r.id - PrefixElementA)})
	if err != nil {
		return err
	}
	_, err = writeString(e, element.Name.Local)
	return err
}

// 0x41
type elementRecord struct {
	elementRecordBase
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

// 0x02
type commentRecord struct {
	textRecordBase
}

func (r *commentRecord) getName() string {
	return "commentRecord (0x02)"
}

func (r *commentRecord) decodeText(d *decoder, trd textRecordDecoder) (string, error) {
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

// 0x03
type arrayRecord struct {
	elementRecordBase
}

func (r *arrayRecord) getName() string {
	return "arrayRecord (0x03)"
}

func (r *arrayRecord) decodeElement(d *decoder) (record, error) {
	rec, err := getNextRecord(d)
	if err != nil {
		return rec, err
	}
	if !rec.isStartElement() {
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
		_, err = valDecoder.decodeText(d, valDecoder)
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
