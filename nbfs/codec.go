package nbfs

import "github.com/khoad/msbingo/nbfx"

// NewDecoder creates a new NBFS Decoder
func NewDecoder() nbfx.Decoder {
	return nbfx.NewDecoderWithStrings(nbfsDictionary)
}

// NewEncoder creates a new NBFS Encoder
func NewEncoder() nbfx.Encoder {
	return nbfx.NewEncoderWithStrings(nbfsDictionary)
}
