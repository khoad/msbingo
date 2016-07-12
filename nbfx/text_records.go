package nbfx

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
	return nil
}

type trueTextRecord struct {
	textRecordBase
}

func (r *trueTextRecord) readText(d *decoder) (string, error) {
	return "true", nil
}

func (r *trueTextRecord) writeText(e *encoder, text string) error {
	return nil
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
