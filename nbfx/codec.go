package nbfx

type codec struct {
	dict map[uint32]string
}

type Encoder interface {
	Encode(xml string) ([]byte, error)
}

type Decoder interface {
	Decode(bin []byte) (string, error)
}

func (c *codec) addDictionaryString(index uint32, value string) {
	c.dict[index] = value
}
