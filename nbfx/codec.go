package nbfx

type Encoder interface {
	Encode(xml string) []byte
}

type encoder struct {
}

func NewEncoder() Encoder {
	return &encoder{}
}

func (e *encoder) Encode(xml string) []byte {
	return []byte{}
}


type Decoder interface {
	Decode(bin []byte) string
}

type decoder struct {
}

func NewDecoder() Decoder {
	return &decoder{}
}

func (d *decoder) Decode(bin []byte) string {
		return ""
}