package nbfx

import (
	"encoding/xml"
	"errors"
	"fmt"
)

func init() {
	records[zeroText] = &zeroTextRecord{textRecordBase{recordBase{"ZeroText", zeroText}, false}}
	records[zeroTextWithEndElement] = &zeroTextRecord{textRecordBase{recordBase{"ZeroTextWithEndElement", zeroTextWithEndElement}, true}}
	records[oneText] = &oneTextRecord{textRecordBase{recordBase{"OneText", oneText}, false}}
	records[oneTextWithEndElement] = &oneTextRecord{textRecordBase{recordBase{"OneTextWithEndElement", oneTextWithEndElement}, true}}
	records[falseText] = &falseTextRecord{textRecordBase{recordBase{"FalseText", falseText}, false}}
	records[falseTextWithEndElement] = &falseTextRecord{textRecordBase{recordBase{"FalseTextWithEndElement", falseTextWithEndElement}, true}}
	records[trueText] = &trueTextRecord{textRecordBase{recordBase{"TrueText", trueText}, false}}
	records[trueTextWithEndElement] = &trueTextRecord{textRecordBase{recordBase{"TrueTextWithEndElement", trueTextWithEndElement}, true}}
	records[int8Text] = &int8TextRecord{textRecordBase{recordBase{"Int8Text", int8Text}, false}}
	records[int8TextWithEndElement] = &int8TextRecord{textRecordBase{recordBase{"Int8TextWithEndElement", int8TextWithEndElement}, true}}
	records[int16Text] = &int16TextRecord{textRecordBase{recordBase{"Int16Text", int16Text}, false}}
	records[int16TextWithEndElement] = &int16TextRecord{textRecordBase{recordBase{"Int16TextWithEndElement", int16TextWithEndElement}, true}}
	records[int32Text] = &int32TextRecord{textRecordBase{recordBase{"Int32Text", int32Text}, false}}
	records[int32TextWithEndElement] = &int32TextRecord{textRecordBase{recordBase{"Int32TextWithEndElement", int32TextWithEndElement}, true}}
	records[int64Text] = &int64TextRecord{textRecordBase{recordBase{"Int64Text", int64Text}, false}}
	records[int64TextWithEndElement] = &int64TextRecord{textRecordBase{recordBase{"Int64TextWithEndElement", int64TextWithEndElement}, true}}
	records[floatText] = &floatTextRecord{textRecordBase{recordBase{"FloatText", floatText}, false}}
	records[floatTextWithEndElement] = &floatTextRecord{textRecordBase{recordBase{"FloatTextWithEndElement", floatTextWithEndElement}, true}}
	records[doubleText] = &doubleTextRecord{textRecordBase{recordBase{"DoubleText", doubleText}, false}}
	records[doubleTextWithEndElement] = &doubleTextRecord{textRecordBase{recordBase{"DoubleTextWithEndElement", doubleTextWithEndElement}, true}}
	records[decimalText] = &decimalTextRecord{textRecordBase{recordBase{"DecimalText", decimalText}, false}}
	records[decimalTextWithEndElement] = &decimalTextRecord{textRecordBase{recordBase{"DecimalTextWithEndElement", decimalTextWithEndElement}, true}}
	records[dateTimeText] = &dateTimeTextRecord{textRecordBase{recordBase{"DateTimeText", dateTimeText}, false}}
	records[dateTimeTextWithEndElement] = &dateTimeTextRecord{textRecordBase{recordBase{"DateTimeTextWithEndElement", dateTimeTextWithEndElement}, true}}
	records[chars8Text] = &chars8TextRecord{textRecordBase{recordBase{"Chars8Text", chars8Text}, false}}
	records[chars8TextWithEndElement] = &chars8TextRecord{textRecordBase{recordBase{"Chars8TextWithEndElement", chars8TextWithEndElement}, true}}
	records[chars16Text] = &chars16TextRecord{textRecordBase{recordBase{"Chars16Text", chars16Text}, false}}
	records[chars16TextWithEndElement] = &chars16TextRecord{textRecordBase{recordBase{"Chars16TextWithEndElement", chars16TextWithEndElement}, true}}
	records[chars32Text] = &chars32TextRecord{textRecordBase{recordBase{"Chars32Text", chars32Text}, false}}
	records[chars32TextWithEndElement] = &chars32TextRecord{textRecordBase{recordBase{"Chars32TextWithEndElement", chars32TextWithEndElement}, true}}
	records[bytes8Text] = &bytes8TextRecord{textRecordBase{recordBase{"Bytes8Text", bytes8Text}, false}}
	records[bytes8TextWithEndElement] = &bytes8TextRecord{textRecordBase{recordBase{"Bytes8TextWithEndElement", bytes8TextWithEndElement}, true}}
	records[bytes16Text] = &bytes16TextRecord{textRecordBase{recordBase{"Bytes16Text", bytes16Text}, false}}
	records[bytes16TextWithEndElement] = &bytes16TextRecord{textRecordBase{recordBase{"Bytes16TextWithEndElement", bytes16TextWithEndElement}, true}}
	records[bytes32Text] = &bytes32TextRecord{textRecordBase{recordBase{"Bytes32Text", bytes32Text}, false}}
	records[bytes32TextWithEndElement] = &bytes32TextRecord{textRecordBase{recordBase{"Bytes32TextWithEndElement", bytes32TextWithEndElement}, true}}
	records[startListText] = &startListTextRecord{textRecordBase{recordBase{"StartListText", startListText}, false}}
	records[startListTextWithEndElement] = &startListTextRecord{textRecordBase{recordBase{"StartListTextWithEndElement", startListTextWithEndElement}, true}}
	records[endListText] = &endListTextRecord{textRecordBase{recordBase{"EndListText", endListText}, false}}
	records[endListTextWithEndElement] = &endListTextRecord{textRecordBase{recordBase{"EndListTextWithEndElement", endListTextWithEndElement}, true}}
	records[emptyText] = &emptyTextRecord{textRecordBase{recordBase{"EmptyText", emptyText}, false}}
	records[emptyTextWithEndElement] = &emptyTextRecord{textRecordBase{recordBase{"EmptyTextWithEndElement", emptyTextWithEndElement}, true}}
	records[dictionaryText] = &dictionaryTextRecord{textRecordBase{recordBase{"DictionaryText", dictionaryText}, false}}
	records[dictionaryTextWithEndElement] = &dictionaryTextRecord{textRecordBase{recordBase{"DictionaryTextWithEndElement", dictionaryTextWithEndElement}, true}}
	records[uniqueIdText] = &uniqueIdTextRecord{textRecordBase{recordBase{"UniqueIdText", uniqueIdText}, false}}
	records[uniqueIdTextWithEndElement] = &uniqueIdTextRecord{textRecordBase{recordBase{"UniqueIdTextWithEndElement", uniqueIdTextWithEndElement}, true}}
	records[timeSpanText] = &timeSpanTextRecord{textRecordBase{recordBase{"TimeSpanText", timeSpanText}, false}}
	records[timeSpanTextWithEndElement] = &timeSpanTextRecord{textRecordBase{recordBase{"TimeSpanTextWithEndElement", timeSpanTextWithEndElement}, true}}
	records[uuidText] = &uuidTextRecord{textRecordBase{recordBase{"UuidText", uuidText}, false}}
	records[uuidTextWithEndElement] = &uuidTextRecord{textRecordBase{recordBase{"UuidTextWithEndElement", uuidTextWithEndElement}, true}}
	records[uInt64Text] = &uInt64TextRecord{textRecordBase{recordBase{"UInt64Text", uInt64Text}, false}}
	records[uInt64TextWithEndElement] = &uInt64TextRecord{textRecordBase{recordBase{"UInt64TextWithEndElement", uInt64TextWithEndElement}, true}}
	records[boolText] = &boolTextRecord{textRecordBase{recordBase{"BoolText", boolText}, false}}
	records[boolTextWithEndElement] = &boolTextRecord{textRecordBase{recordBase{"BoolTextWithEndElement", boolTextWithEndElement}, true}}
	records[unicodeChars8Text] = &unicodeChars8TextRecord{textRecordBase{recordBase{"UnicodeChars8Text", unicodeChars8Text}, false}}
	records[unicodeChars8TextWithEndElement] = &unicodeChars8TextRecord{textRecordBase{recordBase{"UnicodeChars8TextWithEndElement", unicodeChars8TextWithEndElement}, true}}
	records[unicodeChars16Text] = &unicodeChars16TextRecord{textRecordBase{recordBase{"UnicodeChars16Text", unicodeChars16Text}, false}}
	records[unicodeChars16TextWithEndElement] = &unicodeChars16TextRecord{textRecordBase{recordBase{"UnicodeChars16TextWithEndElement", unicodeChars16TextWithEndElement}, true}}
	records[unicodeChars32Text] = &unicodeChars32TextRecord{textRecordBase{recordBase{"UnicodeChars32Text", unicodeChars32Text}, false}}
	records[unicodeChars32TextWithEndElement] = &unicodeChars32TextRecord{textRecordBase{recordBase{"UnicodeChars32TextWithEndElement", unicodeChars32TextWithEndElement}, true}}
	records[qNameDictionaryText] = &qNameDictionaryTextRecord{textRecordBase{recordBase{"QNameDictionaryText", qNameDictionaryText}, false}}
	records[qNameDictionaryTextWithEndElement] = &qNameDictionaryTextRecord{textRecordBase{recordBase{"QNameDictionaryTextWithEndElement", qNameDictionaryTextWithEndElement}, true}}
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
		rec, err := getRecord(endElement)
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

type zeroTextRecord struct {
	textRecordBase
}

func (r *zeroTextRecord) readText(d *decoder) (string, error) {
	return "0", nil
}

func (r *zeroTextRecord) writeText(e *encoder, text string) error {
	return nil
}

type oneTextRecord struct {
	textRecordBase
}

func (r *oneTextRecord) readText(d *decoder) (string, error) {
	return "1", nil
}

func (r *oneTextRecord) writeText(e *encoder, text string) error {
	return nil
}

type falseTextRecord struct {
	textRecordBase
}

func (r *falseTextRecord) readText(d *decoder) (string, error) {
	return "false", nil
}

func (r *falseTextRecord) writeText(e *encoder, text string) error {
	return e.bin.WriteByte(falseText)
}

type trueTextRecord struct {
	textRecordBase
}

func (r *trueTextRecord) readText(d *decoder) (string, error) {
	return "true", nil
}

func (r *trueTextRecord) writeText(e *encoder, text string) error {
	return e.bin.WriteByte(trueText)
}

type int8TextRecord struct {
	textRecordBase
}

func (r *int8TextRecord) readText(d *decoder) (string, error) {
	return readInt8Text(d)
}

type int16TextRecord struct {
	textRecordBase
}

func (r *int16TextRecord) readText(d *decoder) (string, error) {
	return readInt16Text(d)
}

type int32TextRecord struct {
	textRecordBase
}

func (r *int32TextRecord) readText(d *decoder) (string, error) {
	return readInt32Text(d)
}

type int64TextRecord struct {
	textRecordBase
}

func (r *int64TextRecord) readText(d *decoder) (string, error) {
	return readInt64Text(d)
}

type floatTextRecord struct {
	textRecordBase
}

func (r *floatTextRecord) readText(d *decoder) (string, error) {
	return readFloatText(d)
}

type doubleTextRecord struct {
	textRecordBase
}

func (r *doubleTextRecord) readText(d *decoder) (string, error) {
	return readDoubleText(d)
}

type decimalTextRecord struct {
	textRecordBase
}

func (r *decimalTextRecord) readText(d *decoder) (string, error) {
	return readDecimalText(d)
}

type dateTimeTextRecord struct {
	textRecordBase
}

func (r *dateTimeTextRecord) readText(d *decoder) (string, error) {
	return readDateTimeText(d)
}

type chars8TextRecord struct {
	textRecordBase
}

func (r *chars8TextRecord) readText(d *decoder) (string, error) {
	return readChars8Text(d)
}

func (r *chars8TextRecord) writeText(e *encoder, text string) error {
	return writeChars8Text(e, text)
}

type chars16TextRecord struct {
	textRecordBase
}

func (r *chars16TextRecord) readText(d *decoder) (string, error) {
	return readChars16Text(d)
}

type chars32TextRecord struct {
	textRecordBase
}

func (r *chars32TextRecord) readText(d *decoder) (string, error) {
	return readChars32Text(d)
}

func (r *chars32TextRecord) writeText(e *encoder, text string) error {
	return writeChars32Text(e, text)
}

type bytes8TextRecord struct {
	textRecordBase
}

func (r *bytes8TextRecord) readText(d *decoder) (string, error) {
	return readBytes8Text(d)
}

type bytes16TextRecord struct {
	textRecordBase
}

func (r *bytes16TextRecord) readText(d *decoder) (string, error) {
	return readBytes16Text(d)
}

type bytes32TextRecord struct {
	textRecordBase
}

func (r *bytes32TextRecord) readText(d *decoder) (string, error) {
	return readBytes32Text(d)
}

type startListTextRecord struct {
	textRecordBase
}

func (r *startListTextRecord) readText(d *decoder) (string, error) {
	return readListText(d)
}

type endListTextRecord struct {
	textRecordBase
}

type emptyTextRecord struct {
	textRecordBase
}

func (r *emptyTextRecord) readText(d *decoder) (string, error) {
	return "", nil
}

func (r *emptyTextRecord) writeText(e *encoder, text string) error {
	return nil
}

type dictionaryTextRecord struct {
	textRecordBase
}

func (r *dictionaryTextRecord) readText(d *decoder) (string, error) {
	return readDictionaryString(d)
}

func (r *dictionaryTextRecord) writeText(e *encoder, text string) error {
	_, err := writeDictionaryString(e, text)
	return err
}

type uniqueIdTextRecord struct {
	textRecordBase
}

func (r *uniqueIdTextRecord) readText(d *decoder) (string, error) {
	return readUniqueIdText(d)
}

func (r *uniqueIdTextRecord) writeText(e *encoder, text string) error {
	return writeUniqueIdText(e, text)
}

type timeSpanTextRecord struct {
	textRecordBase
}

func (r *timeSpanTextRecord) readText(d *decoder) (string, error) {
	return readTimeSpanText(d)
}

type uuidTextRecord struct {
	textRecordBase
}

func (r *uuidTextRecord) readText(d *decoder) (string, error) {
	return readUuidText(d)
}

func (r *uuidTextRecord) writeText(e *encoder, text string) error {
	return writeUuidText(e, text)
}

type uInt64TextRecord struct {
	textRecordBase
}

func (r *uInt64TextRecord) readText(d *decoder) (string, error) {
	return readUInt64Text(d)
}

type boolTextRecord struct {
	textRecordBase
}

func (r *boolTextRecord) readText(d *decoder) (string, error) {
	return readBoolText(d)
}

func (r *boolTextRecord) writeText(e *encoder, text string) error {
	//if text == "false" {
	//	return e.bin.WriteByte(0)
	//} else if text == "true" {
	//	return e.bin.WriteByte(1)
	//}
	//return errors.New("BoolText record text must be 'true' or 'false'")
	return errors.New("boolTextRecord.writeText: Not Implemented")
}

type unicodeChars8TextRecord struct {
	textRecordBase
}

func (r *unicodeChars8TextRecord) readText(d *decoder) (string, error) {
	return readUnicodeChars8Text(d)
}

type unicodeChars16TextRecord struct {
	textRecordBase
}

func (r *unicodeChars16TextRecord) readText(d *decoder) (string, error) {
	return readUnicodeChars16Text(d)
}

type unicodeChars32TextRecord struct {
	textRecordBase
}

func (r *unicodeChars32TextRecord) readText(d *decoder) (string, error) {
	return readUnicodeChars32Text(d)
}

type qNameDictionaryTextRecord struct {
	textRecordBase
}

func (r *qNameDictionaryTextRecord) readText(d *decoder) (string, error) {
	return readQNameDictionaryText(d)
}
