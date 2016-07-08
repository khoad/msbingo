package nbfx

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
)

type encoder struct {
	dict        map[string]uint32
	xml *xml.Decoder
	bin *io.Writer
}

func (e *encoder) addDictionaryString(index uint32, value string) {
	if _, ok := e.dict[value]; ok {
		return
	}
	e.dict[value] = index
}

func NewEncoder() Encoder {
	return NewEncoderWithStrings(nil)
}

func NewEncoderWithStrings(dictionaryStrings map[uint32]string) Encoder {
	encoder := &encoder{make(map[string]uint32), nil, nil}
	if dictionaryStrings != nil {
		for k, v := range dictionaryStrings {
			encoder.addDictionaryString(k, v)
		}
	}
	return encoder
}

func (e *encoder) Encode(xmlString string) ([]byte, error) {
	bin := &bytes.Buffer{}
	reader := bytes.NewReader([]byte(xmlString))
	xml := xml.NewDecoder(reader)
	token, err := xml.RawToken()
	for err == nil {
		record, err := e.getRecordFromToken(token)
		if err != nil {
			return bin.Bytes(), err
		}
		fmt.Println("Encode record", record.getName())
		if record.isElement() {
			elementWriter := record.(elementRecordWriter)
			err = elementWriter.writeElement(e)
		} else {
			textWriter := record.(textRecordWriter)
			err = textWriter.writeText(e)
		}
		if err != nil {
			return bin.Bytes(), errors.New(fmt.Sprintf("Error writing Token %s :: %s", token, err.Error()))
		}
		token, err = xml.RawToken()
	}
	return bin.Bytes(), nil
}

func (e *encoder) getRecordFromToken(token xml.Token) (record, error) {
	switch token.(type) {
	case xml.StartElement:
		return e.getStartElementRecordFromToken(token.(xml.StartElement))
	case xml.CharData:
		return e.getTextRecordFromToken(token.(xml.CharData))
	}

	tokenXmlBytes, err := xml.Marshal(token)
	var tokenXml string
	if err != nil {
		tokenXml = "[[UNKNOWN]]"
	} else {
		tokenXml = string(tokenXmlBytes)
	}
	return nil, errors.New(fmt.Sprint("Unknown token", tokenXml))
}

func (e *encoder) getTextRecordFromToken(cd xml.CharData) (record, error) {
	//return records[0x9C](c)
	return nil, errors.New("UnsupportedOpertation: getTextRecordFromToken")
}

func (e *encoder) getStartElementRecordFromToken(startElement xml.StartElement) (record, error) {
	//fmt.Printf("Getting start element for %s", startElement.Name.Local)
	prefix := startElement.Name.Space
	name := startElement.Name.Local
	prefixIndex := -1
	if len(prefix) == 1 && byte(prefix[0]) >= byte('a') && byte(prefix[0]) <= byte('z') {
		prefixIndex = int(byte(prefix[0]) - byte('a'))
	}
	isNameIndexAssigned := false
	if _, ok := e.dict[name]; ok {
		isNameIndexAssigned = true
	}

	if prefix == "" {
		if !isNameIndexAssigned {
			return records[ShortElement], nil
		} else {
			return records[ShortDictionaryElement], nil
		}
	} else if prefixIndex != -1 {
		if !isNameIndexAssigned {
			return records[PrefixElementA + byte(prefixIndex)], nil
		} else {
			return records[PrefixDictionaryElementA+ byte(prefixIndex)], nil
		}
	} else {
		if !isNameIndexAssigned {
			return records[Element], nil
		} else {
			return records[DictionaryElement], nil
		}
	}
	return nil, errors.New("getStartElementRecordFromToken unable to resolve required xml.Token")
}

func writeString(writer io.Writer, str string) (int, error) {
	var strBytes = []byte(str)
	lenByteLen, err := writeMultiByteInt31(writer, uint32(len(strBytes)))
	if err != nil {
		return lenByteLen, err
	}
	strByteLen, err := writer.Write(strBytes)
	return lenByteLen + strByteLen, err
}

func writeMultiByteInt31(writer io.Writer, num uint32) (int, error) {
	max := uint32(2147483647)
	if num > max {
		return 0, errors.New(fmt.Sprintf("Overflow: i (%d) must be <= max (%d)", num, max))
	}
	if num < MASK_MBI31 {
		return writer.Write([]byte{byte(num)})
	}
	q := num / MASK_MBI31
	rem := num % MASK_MBI31
	n1, err := writer.Write([]byte{byte(MASK_MBI31 + rem)})
	if err != nil {
		return n1, err
	}
	n2, err := writeMultiByteInt31(writer, q)
	return n1 + n2, err
}
