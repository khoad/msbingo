package nbfx

import (
	"encoding/xml"
	"errors"
	"fmt"
)

func init() {
	records[shortElement] = &shortElementRecord{elementRecordBase{recordBase{"ShortElement", shortElement}}}
	records[element] = &elementRecord{elementRecordBase{recordBase{"Element", element}}}
	records[shortDictionaryElement] = &shortDictionaryElementRecord{elementRecordBase{recordBase{"ShortDictionaryElement", shortDictionaryElement}}}
	records[dictionaryElement] = &dictionaryElementRecord{elementRecordBase{recordBase{"DictionaryElement", dictionaryElement}}}
	addAzRecords(prefixDictionaryElementA, "PrefixDictionaryElement", func(id byte, name string) record {
		return &prefixDictionaryElementAZRecord{elementRecordBase{recordBase{name, id}}}
	})
	addAzRecords(prefixElementA, "PrefixElement", func(id byte, name string) record {
		return &prefixElementAZRecord{elementRecordBase{recordBase{name, id}}}
	})
}

type elementRecordBase struct {
	recordBase
}

func (r *elementRecordBase) isStartElement() bool { return true }

func (r *elementRecordBase) readElementAttributes(element xml.StartElement, d *decoder) (record, error) {
	var peekRecord record
	var attributeToken xml.Attr

	// get next record
	rec, err := getNextRecord(d)
	for err == nil && rec != nil {
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

	d.elementStack.push(element)

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

//(0x44-0x5D)
type prefixDictionaryElementAZRecord struct {
	elementRecordBase
}

func (r *prefixDictionaryElementAZRecord) decodeElement(d *decoder) (record, error) {
	name, err := readDictionaryString(d)
	if err != nil {
		return nil, err
	}
	element := xml.StartElement{Name: xml.Name{Local: string('a'+r.id-prefixDictionaryElementA) + ":" + name}}

	return r.readElementAttributes(element, d)
}

func (r *prefixDictionaryElementAZRecord) encodeElement(e *encoder, element xml.StartElement) error {
	_, err := e.bin.Write([]byte{r.id})
	if err != nil {
		return err
	}
	err = writeDictionaryString(e, element.Name.Local)
	if err != nil {
		return err
	}
	return r.encodeAttributes(e, element.Attr)
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

func (r *shortDictionaryElementRecord) encodeElement(e *encoder, element xml.StartElement) error {
	_, err := e.bin.Write([]byte{r.id})
	if err != nil {
		return err
	}
	err = writeDictionaryString(e, element.Name.Local)
	if err != nil {
		return err
	}
	return r.encodeAttributes(e, element.Attr)
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

func (r *dictionaryElementRecord) encodeElement(e *encoder, element xml.StartElement) error {
	err := e.bin.WriteByte(r.id)
	if err != nil {
		return err
	}
	_, err = writeString(e, element.Name.Space)
	if err != nil {
		return err
	}
	err = writeDictionaryString(e, element.Name.Local)
	if err != nil {
		return err
	}
	return r.encodeAttributes(e, element.Attr)
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
	err := e.bin.WriteByte(r.id)
	if err != nil {
		return err
	}
	_, err = writeString(e, element.Name.Local)
	if err != nil {
		return err
	}
	return r.encodeAttributes(e, element.Attr)
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
	element := xml.StartElement{Name: xml.Name{Local: string('a'+(r.id-prefixElementA)) + ":" + name}}

	return r.readElementAttributes(element, d)
}

func (r *prefixElementAZRecord) encodeElement(e *encoder, element xml.StartElement) error {
	err := e.bin.WriteByte(prefixElementA + (r.id - prefixElementA))
	if err != nil {
		return err
	}
	_, err = writeString(e, element.Name.Local)
	if err != nil {
		return err
	}
	return r.encodeAttributes(e, element.Attr)
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

func (r *elementRecord) encodeElement(e *encoder, element xml.StartElement) error {
	err := e.bin.WriteByte(r.id)
	if err != nil {
		return err
	}
	_, err = writeString(e, element.Name.Space)
	if err != nil {
		return err
	}
	_, err = writeString(e, element.Name.Local)
	if err != nil {
		return err
	}
	return r.encodeAttributes(e, element.Attr)
}
