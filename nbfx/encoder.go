package nbfx

import (
	"errors"
)

type encoder struct {
	codec codec
}

func (e *encoder) Encode(xml string) ([]byte, error) {
	return []byte{}, errors.New("NotImplemented: nbfx.Encoder.Encode(string)")
}
