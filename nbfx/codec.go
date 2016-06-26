package nbfx

type codec struct {
}

type Encoder interface {
	Encode(xml string) ([]byte, error)
}

type Decoder interface {
	Decode(bin []byte) (string, error)
}

var dict = map[uint32]string{
	2: "Envelope",
	8: "Header",
}

func addDictionaryString(index uint32, value string) {
	dict[index] = value
}
