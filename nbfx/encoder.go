package nbfx

import (
	"bytes"
	"errors"
	"encoding/xml"
	"fmt"
	"strings"
)

type encoder struct {
	codec codec
}

func NewEncoder() Encoder {
	return NewEncoderWithStrings(nil)
}

func NewEncoderWithStrings(dictionaryStrings map[uint32]string) Encoder {
	encoder := &encoder{codec{make(map[uint32]string)}}
	if dictionaryStrings != nil {
		for k, v := range dictionaryStrings {
			encoder.codec.addDictionaryString(k, v)
		}
	}
	return encoder
}

func (e *encoder) Encode(xmlString string) ([]byte, error) {
	reader := bytes.NewReader([]byte(xmlString))
	binBuffer := &bytes.Buffer{}
	xmlDecoder := xml.NewDecoder(reader)
	token, err := xmlDecoder.RawToken()
	for err == nil {
		record := getRecordFromToken(encoder.codec, token)
		if record == nil {
			return binBuffer.Bytes(), errors.New(fmt.Sprintf("Unknown Token %s", token))
		}
	}
	return []byte{}, errors.New("NotImplemented: nbfx.Encoder.Encode(string)")
}

func getRecordFromToken(codec codec, token xml.Token) record {
	switch token.(type) {
		case xml.StartElement:
			return getStartElementRecordFromToken(codec, token.(xml.StartElement))
	}

	return nil
}

func getStartElementRecordFromToken(codec *codec, startElement xml.StartElement) record {
	name := startElement.Name.Local
	parts := strings.Split(name, ":")
	prefix := ""
	if len(parts) > 1 {
		prefix = parts[0]
		name = parts[1]
	}
	isPrefixAZ := -1
	if len(prefix) == 1 && byte(prefix) >= byte('a') && byte(prefix) <= byte('z') {
		isPrefixAZ = byte(prefix) - byte('a')
	}
	nameIndex := -1
	if i, ok := codec.reverseDict[name]; ok {
		nameIndex = i
	}
	if nameIndex >= 0 {
		if prefix == "" {
			return dictionaryElement{nameIndex: nameIndex, name: name}
		} else if (isPrefixAZ >= 0){
			return prefixDictionaryElementS{prefixIndex: isPrefixAZ, prefix: prefix, nameIndex: nameIndex, name: name}
		} else {
			return dictionaryElement{prefix: prefix, nameIndex: nameIndex, name: name}
		}
	}
	if prefix == "" {
		return shortElement{name: name}
	} else {
		return prefixShortElement{prefix: prefix}
	}
}
