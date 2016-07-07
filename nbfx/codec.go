package nbfx

type Encoder interface {
	Encode(xml string) ([]byte, error)
}

type Decoder interface {
	Decode(bin []byte) (string, error)
}

// for MultiByteInt31
const MASK_MBI31 = uint32(0x80) //0x80 = 128
