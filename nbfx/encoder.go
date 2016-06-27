package nbfx

import (
	"bytes"
	"errors"
	"encoding/xml"
	"fmt"
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
		record := getRecordFromToken(token)
		if record == nil {
			return binBuffer.Bytes(), errors.New(fmt.Sprintf("Unknown Token %s", token))
		}
	}
	return []byte{}, errors.New("NotImplemented: nbfx.Encoder.Encode(string)")
}

func getRecordFromToken(token xml.Token) record {
	switch token.(type) {
		case xml.StartElement:
			return getStartElementRecordFromToken(token)
	}

	return nil
}

func getStartElementRecordFromToken(token xml.Token) record {
	return nil
}
