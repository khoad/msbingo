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

type elementRecordWriter interface {
	encodeElement(e *encoder, element xml.StartElement) error
}

type attributeRecordDecoder interface {
	decodeAttribute(d *decoder) (xml.Attr, error)
}

type textRecordDecoder interface {
	decodeText(d *decoder) (string, error)
	readText(d *decoder) (string, error)
}

type textRecordWriter interface {
	encodeText(e *encoder, cd xml.CharData) error
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
	//fmt.Println("got next record", peekRecord, err)

	err = d.xml.EncodeToken(element)
	if err != nil {
		return nil, err
	}

	d.elementStack.Push(element)

	return peekRecord, nil
}

func (r *elementRecordBase) encodeElement(e *encoder, element xml.StartElement) error {
	return errors.New(fmt.Sprint("NotImplemented: encodeElement on", r))
}

type attributeRecordBase struct{
	recordBase
}

func (r *attributeRecordBase) isAttribute() bool { return true }

type textRecordBase struct {
	recordBase
	withEndElement bool
	textReader       func(d *decoder) (string, error)
	textWriter func(e *encoder, text string) error
}

func (r *textRecordBase) isText() bool {
	return true
}

func (r *textRecordBase) getName() string {
	return fmt.Sprintf("%s (%#x)", r.name, r.id)
}

func (r *textRecordBase) readText(d *decoder) (string, error) {
	if r.textReader == nil {
		return "", errors.New(fmt.Sprint("Error", r, "textReader is nil"))
	}
	return r.textReader(d)
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

func (r *textRecordBase) encodeText(e *encoder, cd xml.CharData) error {
	if r.textWriter == nil {
		b, err := xml.Marshal(cd)
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		str := string(b)
		return errors.New(fmt.Sprint("NotImplement: writeText for " + r.getName() + " :: [" + str + "]", errMsg))
	}

	return r.textWriter(e, string(cd))
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
	records[Comment] = &commentRecord{textRecordBase{recordBase{"Comment",Comment},false,nil,nil}}
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
	addTextRecord(ZeroText, "ZeroText", func(d *decoder) (string, error) { return "0", nil }, nil)
	addTextRecord(OneText, "OneText", func(d *decoder) (string, error) { return "1", nil }, nil)
	addTextRecord(FalseText, "FalseText", func(d *decoder) (string, error) { return "false", nil }, nil)
	addTextRecord(TrueText, "TrueText", func(d *decoder) (string, error) { return "true", nil}, nil)
	addTextRecord(Int8Text, "Int8Text", func(d *decoder) (string, error) { return readInt8Text(d) }, nil)
	addTextRecord(Int16Text, "Int16Text", func(d *decoder) (string, error) { return readInt16Text(d) }, nil)
	addTextRecord(Int32Text, "Int32Text", func(d *decoder) (string, error) { return readInt32Text(d) }, nil)
	addTextRecord(Int64Text, "Int64Text", func(d *decoder) (string, error) { return readInt64Text(d) }, nil)
	addTextRecord(FloatText, "FloatText", func(d *decoder) (string, error) { return readFloatText(d) }, nil)
	addTextRecord(DoubleText, "DoubleText", func(d *decoder) (string, error) { return readDoubleText(d) }, nil)
	addTextRecord(DecimalText, "DecimalText", func(d *decoder) (string, error) { return readDecimalText(d) }, nil)
	addTextRecord(DateTimeText, "DateTimeText", func(d *decoder) (string, error) { return readDateTimeText(d) }, nil)
	addTextRecord(Chars8Text, "Chars8Text", func(d *decoder) (string, error) { return readChars8Text(d) }, nil)
	addTextRecord(Chars16Text, "Chars16Text", func(d *decoder) (string, error) { return readChars16Text(d) }, nil)
	addTextRecord(Chars32Text, "Chars32Text", func(d *decoder) (string, error) { return readChars32Text(d) }, func(e *encoder, text string) error { return writeChars32Text(e, text) })
	addTextRecord(Bytes8Text, "Bytes8Text", func(d *decoder) (string, error) { return readBytes8Text(d) }, nil)
	addTextRecord(Bytes16Text, "Bytes16Text", func(d *decoder) (string, error) { return readBytes16Text(d) }, nil)
	addTextRecord(Bytes32Text, "Bytes32Text", func(d *decoder) (string, error) { return readBytes32Text(d) }, nil)
	addTextRecord(StartListText, "StartListText", func(d *decoder) (string, error) { return readListText(d) }, nil)
	addTextRecord(EndListText, "EndListText", func(d *decoder) (string, error) { return "", nil }, nil)
	addTextRecord(EmptyText, "EmptyText", func(d *decoder) (string, error) { return "", nil }, nil)
	addTextRecord(DictionaryText, "DictionaryText", func(d *decoder) (string, error) { return readDictionaryString(d) }, nil)
	addTextRecord(UniqueIdText, "UniqueIdText", func(d *decoder) (string, error) { return readUniqueIdText(d) }, nil)
	addTextRecord(TimeSpanText, "TimeSpanText", func(d *decoder) (string, error) { return readTimeSpanText(d) }, nil)
	addTextRecord(UuidText, "UuidText", func(d *decoder) (string, error) { return readUuidText(d) }, nil)
	addTextRecord(UInt64Text, "UInt64Text", func(d *decoder) (string, error) { return readUInt64Text(d) }, nil)
	addTextRecord(BoolText, "BoolText", func(d *decoder) (string, error) { return readBoolText(d) }, nil)
	addTextRecord(UnicodeChars8Text, "UnicodeChars8Text", func(d *decoder) (string, error) { return readUnicodeChars8Text(d) }, nil)
	addTextRecord(UnicodeChars16Text, "UnicodeChars16Text", func(d *decoder) (string, error) { return readUnicodeChars16Text(d) }, nil)
	addTextRecord(UnicodeChars32Text, "UnicodeChars32Text", func(d *decoder) (string, error) { return readUnicodeChars32Text(d) }, nil)
	addTextRecord(QNameDictionaryText, "QNameDictionaryText", func(d *decoder) (string, error) { return readQNameDictionaryText(d) }, nil)

	addAzRecords()
}

func addTextRecord(recordId byte, textName string, textReader func(*decoder) (string, error), textWriter func(*encoder,string) error) {
	records[recordId] = &textRecordBase{recordBase{name: textName, id: recordId}, false, textReader, textWriter}
	records[recordId+1] = &textRecordBase{recordBase{name: textName + "WithEndElement", id: recordId + 1}, true, textReader, textWriter}
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
	//fmt.Println(fmt.Sprint("Writing PrefixDictionaryElement" + string('a' + PrefixDictionaryElementA), element))
	e.bin.Write([]byte{r.id})
	_, err := writeMultiByteInt31(e, e.dict[element.Name.Local])
	//if err != nil {
	//	fmt.Println("Write PrefixDictionaryElement error", err.Error())
	//}
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
