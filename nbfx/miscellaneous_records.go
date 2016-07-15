package nbfx

import (
	"encoding/xml"
	"errors"
)

func init() {
	records[endElement] = &endElementRecord{elementRecordBase{recordBase{"EndElement", endElement}}}
	records[comment] = &commentRecord{textRecordBase{recordBase{"Comment", comment}, false}}
	records[array] = &arrayRecord{elementRecordBase{recordBase{"Array", array}}}
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
	endElementToken := xml.EndElement{Name: xml.Name{Local: element.Name.Local, Space: element.Name.Space}}
	err := d.xml.EncodeToken(endElementToken)
	return nil, err
}

func (r *endElementRecord) encodeElement(e *encoder, element xml.StartElement) error {
	_, err := e.bin.Write([]byte{r.id})
	return err
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
		if i == 0 {
			startElement = d.elementStack.top.value.(xml.StartElement)
		} else {
			err = d.xml.EncodeToken(startElement)
			if err != nil {
				return nil, err
			}
			d.elementStack.Push(startElement)
		}
		_, err = valDecoder.decodeText(d, valDecoder)
		if err != nil {
			return nil, err
		}
		if i < len {
			d.elementStack.Push(startElement)
		}
	}
	return nil, nil
}
