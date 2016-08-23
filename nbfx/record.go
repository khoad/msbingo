package nbfx

import (
	"encoding/xml"
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
	name string
	id   byte
}

func (r *recordBase) isStartElement() bool { return false }
func (r *recordBase) isEndElement() bool   { return false }
func (r *recordBase) isAttribute() bool    { return false }
func (r *recordBase) isText() bool         { return false }
func (r *recordBase) getName() string      { return r.name }

func getNextRecord(d *decoder) (record, error) {
	b, err := readByte(d.bin)
	if err != nil {
		return nil, err
	}

	return getRecord(b)
}

func getRecord(b byte) (record, error) {
	if rec, ok := records[b]; ok {
		return rec, nil
	}

	return nil, fmt.Errorf("Unknown record %#x", b)
}

var records = make(map[byte]record)

func addAzRecords(idA byte, baseName string, recFunc func(byte, string) record) {
	for i := 0; i < 26; i++ {
		id := idA + byte(i)
		rec := recFunc(id, baseName+strings.ToUpper(string('a'+i)))
		records[id] = rec
	}
}
