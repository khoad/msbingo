package nbfx

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func init() {
	records[shortAttribute] = &shortAttributeRecord{attributeRecordBase{recordBase{"ShortAttribute", shortAttribute}}}
	records[attribute] = &attributeRecord{attributeRecordBase{recordBase{"Attribute", attribute}}}
	records[shortDictionaryAttribute] = &shortDictionaryAttributeRecord{attributeRecordBase{recordBase{"ShortDictionaryAttribute", shortDictionaryAttribute}}}
	records[dictionaryAttribute] = &dictionaryAttributeRecord{attributeRecordBase{recordBase{"DictionaryAttribute", dictionaryAttribute}}}
	records[shortXmlnsAttribute] = &shortXmlnsAttributeRecord{attributeRecordBase{recordBase{"ShortXmlnsAttribute", shortXmlnsAttribute}}}
	records[xmlnsAttribute] = &xmlnsAttributeRecord{attributeRecordBase{recordBase{"XmlnsAttribute", xmlnsAttribute}}}
	records[shortDictionaryXmlnsAttribute] = &shortDictionaryXmlnsAttributeRecord{attributeRecordBase{recordBase{"ShortDictionaryXmlnsAttribute", shortDictionaryXmlnsAttribute}}}
	records[dictionaryXmlnsAttribute] = &dictionaryXmlnsAttributeRecord{attributeRecordBase{recordBase{"DictionaryXmlnsAttribute", dictionaryXmlnsAttribute}}}
	addAzRecords(prefixDictionaryAttributeA, "PrefixDictionaryAttribute", func(id byte, name string) record {
		return &prefixDictionaryAttributeAZRecord{attributeRecordBase{recordBase{name, id}}}
	})
	addAzRecords(prefixAttributeA, "PrefixAttribute", func(id byte, name string) record {
		return &prefixAttributeAZRecord{attributeRecordBase{recordBase{name, id}}}
	})
}

type attributeRecordBase struct {
	recordBase
}

func (r *attributeRecordBase) isAttribute() bool { return true }

func (r *attributeRecordBase) encodeAttribute(e *encoder, attr xml.Attr) error {
	return errors.New(fmt.Sprint("NotImplemented: encodeAttribute on", r))
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

func (r *shortAttributeRecord) encodeAttribute(e *encoder, attr xml.Attr) error {
	var err error
	err = e.bin.WriteByte(r.id)
	if err != nil {
		return err
	}
	_, err = writeString(e, attr.Name.Local)
	if err != nil {
		return err
	}

	rec, err := e.getTextRecordFromText(attr.Value, false)
	if err != nil {
		return err
	}

	textWriter := rec.(textRecordEncoder)
	return textWriter.writeText(e, attr.Value)
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

func (r *shortDictionaryAttributeRecord) encodeAttribute(e *encoder, attr xml.Attr) error {
	var err error
	err = e.bin.WriteByte(r.id)
	if err != nil {
		return err
	}

	if val, ok := e.dict[attr.Name.Local]; ok {
		valString := strconv.Itoa(int(val))
		if err != nil {
			return err
		}
		_, err = writeString(e, valString)
		if err != nil {
			return err
		}
	} else if strings.HasPrefix(attr.Name.Local, "str") {
		// capture "8" in "str8" and write "8"
		numString := attr.Name.Local[3:]
		numInt, err := strconv.Atoi(numString)
		if err != nil {
			return err
		}
		err = e.bin.WriteByte(byte(numInt))
		if err != nil {
			return err
		}
	} else {
		return errors.New("Invalid Operation")
	}

	rec, err := e.getTextRecordFromText(attr.Value, false)
	if err != nil {
		return err
	}

	textWriter := rec.(textRecordEncoder)
	return textWriter.writeText(e, attr.Value)
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

func (r *shortXmlnsAttributeRecord) encodeAttribute(e *encoder, attr xml.Attr) error {
	err := e.bin.WriteByte(r.id)
	if err != nil {
		return err
	}
	_, err = writeString(e, attr.Value)
	return err
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

func (r *xmlnsAttributeRecord) encodeAttribute(e *encoder, attr xml.Attr) error {
	err := e.bin.WriteByte(r.id)
	if err != nil {
		return err
	}
	_, err = writeString(e, attr.Name.Local)
	if err != nil {
		return err
	}
	_, err = writeString(e, attr.Value)
	return err
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
	return xml.Attr{Name: xml.Name{Local: string('a'+r.id-prefixDictionaryAttributeA) + ":" + name}, Value: text}, nil
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
	return xml.Attr{Name: xml.Name{Local: string('a'+r.id-prefixAttributeA) + ":" + name}, Value: text}, nil
}
