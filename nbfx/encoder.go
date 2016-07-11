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
	tokenBuffer *Queue
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
	encoder := &encoder{make(map[string]uint32), nil, nil, &Queue{}}
	if dictionaryStrings != nil {
		for k, v := range dictionaryStrings {
			encoder.addDictionaryString(k, v)
		}
	}
	return encoder
}

func (e *encoder) popToken() (xml.Token, error) {
	var token xml.Token
	var err error
	if e.tokenBuffer.Len() > 0 {
		token = e.tokenBuffer.Dequeue().(xml.Token)
		//fmt.Println("Popped", token)
		return token, err
	}
	token, err = e.xml.RawToken()
	if err == nil {
		//fmt.Println("Popped", token)
	}
	token = xml.CopyToken(token) // make the token immutable (see doc for xml.Decoder.Token())
	return token, err
}

func (e *encoder) pushToken(token xml.Token) {
	e.tokenBuffer.Enqueue(token)
}

func (e *encoder) Encode(xmlString string) ([]byte, error) {
	e.bin = &bytes.Buffer{}
	reader := bytes.NewReader([]byte(xmlString))
	e.xml = xml.NewDecoder(reader)
	token, err := e.popToken()
	for err == nil && token != nil {
		record, err := e.getRecordFromToken(token)
		if err != nil {
			return e.bin.Bytes(), err
		}
		//fmt.Println("Encode record", record.getName())
		if record.isStartElement() {
			//fmt.Println("Encode element", token)
			elementWriter := record.(elementRecordEncoder)
			err = elementWriter.encodeElement(e, token.(xml.StartElement))
		} else if record.isText() {
			textWriter := record.(textRecordEncoder)
			err = textWriter.encodeText(e, textWriter, string(token.(xml.CharData)))
		} else if record.isEndElement() {
			elementWriter := record.(elementRecordEncoder)
			err = elementWriter.encodeElement(e, xml.StartElement{})
		} else {
			err = errors.New(fmt.Sprint("NotSupported: Encoding record", record))
		}
		if err != nil {
			return e.bin.Bytes(), errors.New(fmt.Sprintf("Error writing Token %s :: %s", token, err.Error()))
		}
		token, err = e.popToken()
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
	withEndElement := false
	next, _ := e.popToken()
	switch next.(type) {
	case xml.EndElement:
		withEndElement = true
	}
	//fmt.Println("getTextRecordFromToken", string(cd), next, withEndElement)
	text := string(cd)
	if !withEndElement {
		e.pushToken(next)
	}
	return e.getTextRecordFromText(text, withEndElement)
}

func (e *encoder) getTextRecordFromText(text string, withEndElement bool) (record, error) {
	var id byte
	id = 0x00
	if text == "" {
		id = EmptyText
	} else if text == "0" {
		id = ZeroText
	} else if text == "1" {
		id = OneText
	} else if text == "false" {
		id = FalseText
	} else if text == "true" {
		id = TrueText
	} else {
		if _, ok := e.dict[text]; ok {
			id = DictionaryText
		} else {
			id = Chars8Text
		}
	}
	if id != 0 && withEndElement {
		id += 1
	}
	if rec, ok := records[id]; ok {
		return rec, nil
	}
	return nil, errors.New(fmt.Sprintf("Unknown text record id %#X for %s withEndElement %v", id, text, withEndElement))
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
	return nil, errors.New(fmt.Sprint("getStartElementRecordFromToken unable to resolve", startElement))
}

func (e *encoder) getAttributeRecordFromToken(attr xml.Attr) (record, error) {
	//fmt.Printf("Getting attr element for %s", attr.Name.Local)
	prefix := attr.Name.Space
	name := attr.Name.Local
	isXmlns := prefix == "xmlns" || prefix == "" && name == "xmlns"
	prefixIndex := -1
	if len(prefix) == 1 && byte(prefix[0]) >= byte('a') && byte(prefix[0]) <= byte('z') {
		prefixIndex = int(byte(prefix[0]) - byte('a'))
	}
	isNameIndexAssigned := false
	if _, ok := e.dict[name]; ok {
		isNameIndexAssigned = true
	}

	//fmt.Println("getAttributeRecordFromToken", prefix, name, isXmlns, prefixIndex, isNameIndexAssigned)

	if prefix == "" {
		if isXmlns {
			if _, ok := e.dict[attr.Value]; ok {
				return records [ShortDictionaryXmlnsAttribute], nil
			} else {
				return records [ShortXmlnsAttribute], nil
			}
		} else if isNameIndexAssigned {
			return records[ShortDictionaryAttribute], nil
		} else {
			return records[ShortAttribute], nil
		}
	} else if prefixIndex != -1 {
		if !isNameIndexAssigned {
			return records[PrefixAttributeA + byte(prefixIndex)], nil
		} else {
			return records[PrefixDictionaryAttributeA+ byte(prefixIndex)], nil
		}
	} else {
		if isXmlns {
			if !isNameIndexAssigned {
				return records[XmlnsAttribute], nil
			} else {
				return records[DictionaryXmlnsAttribute], nil
			}
		} else {
			if !isNameIndexAssigned {
				return records[Attribute], nil
			} else {
				return records[DictionaryAttribute], nil
			}
		}
	}
	return nil, errors.New(fmt.Sprint("getAttributeRecordFromToken unable to resolve", attr))
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

func writeDictionaryString(e *encoder, str string) (int, error) {
	key, ok := e.dict[str]
	if !ok {
		return 0, errors.New(fmt.Sprint("Value %s not found in dictionary for DictionaryString record", str))
	}
	return writeMultiByteInt31(e, uint32(key))
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

func writeChars8Text(e *encoder, text string) error {
	bytes := []byte(text)
	writeMultiByteInt31(e, uint32(len(bytes)))
	_, err := e.bin.Write(bytes)
	return err
}

func writeChars32Text(e *encoder, text string) error {
	bytes := []byte(text)
	writeMultiByteInt31(e, uint32(len(bytes)))
	_, err := e.bin.Write(bytes)
	return err
}
