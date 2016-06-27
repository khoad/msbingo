package nbfx

import (
	"errors"
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

func (e *encoder) Encode(xml string) ([]byte, error) {
	return []byte{}, errors.New("NotImplemented: nbfx.Encoder.Encode(string)")
}
