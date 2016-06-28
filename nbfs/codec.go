package nbfs

import "github.com/khoad/msbingo/nbfx"

func NewDecoder() nbfx.Decoder {
	return nbfx.NewDecoderWithStrings(nbfsDictionary)
}
