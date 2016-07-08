package nbfx

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
)

type encoder struct {
	dict        map[string]uint32
	xml *xml.Decoder
	bin *bytes.Buffer
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
	e.bin = &bytes.Buffer{}
	reader := bytes.NewReader([]byte(xmlString))
	e.xml = xml.NewDecoder(reader)
	token, err := e.xml.RawToken()
	for err == nil && token != nil {
		record, err := e.getRecordFromToken(token)
		if err != nil {
			return e.bin.Bytes(), err
		}
		//fmt.Println("Encode record", record.getName())
		if record.isStartElement() {
			fmt.Println("Encode elment", token)
			elementWriter := record.(elementRecordWriter)
			fmt.Println("Writer is", elementWriter)
			err = elementWriter.encodeElement(e, token.(xml.StartElement))
		} else if record.isText() {
			textWriter := record.(textRecordWriter)
			err = textWriter.encodeText(e, token.(xml.CharData))
		} else if record.isEndElement() {
			elementWriter := record.(elementRecordWriter)
			err = elementWriter.encodeElement(e, xml.StartElement{})
		} else {
			err = errors.New(fmt.Sprint("NotSupported: Encoding record", record))
		}
		if err != nil {
			return e.bin.Bytes(), errors.New(fmt.Sprintf("Error writing Token %s :: %s", token, err.Error()))
		}
		token, err = e.xml.RawToken()
	}
	return e.bin.Bytes(), nil
}

func (e *encoder) getRecordFromToken(token xml.Token) (record, error) {
	//fmt.Println("getRecordFromToken", token)
	switch token.(type) {
	case xml.StartElement:
		return e.getStartElementRecordFromToken(token.(xml.StartElement))
	case xml.CharData:
		return e.getTextRecordFromToken(token.(xml.CharData))
	case xml.EndElement:
		return records[EndElement], nil
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
	return records[Chars32Text], nil
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

func writeString(e *encoder, str string) (int, error) {
	var strBytes = []byte(str)
	lenByteLen, err := writeMultiByteInt31(e, uint32(len(strBytes)))
	if err != nil {
		return lenByteLen, err
	}
	strByteLen, err := e.bin.Write(strBytes)
	return lenByteLen + strByteLen, err
}

func writeMultiByteInt31(e *encoder, num uint32) (int, error) {
	max := uint32(2147483647)
	if num > max {
		return 0, errors.New(fmt.Sprintf("Overflow: i (%d) must be <= max (%d)", num, max))
	}
	if num < MASK_MBI31 {
		return e.bin.Write([]byte{byte(num)})
	}
	q := num / MASK_MBI31
	rem := num % MASK_MBI31
	n1, err := e.bin.Write([]byte{byte(MASK_MBI31 + rem)})
	if err != nil {
		return n1, err
	}
	n2, err := writeMultiByteInt31(e, q)
	return n1 + n2, err
}

func writeChars32Text(e *encoder, text string) error {
	_, err := e.bin.Write([]byte{Chars32Text})
	if err != nil {
		return err
	}
	_, err = e.bin.Write([]byte(text))
	return err
}
