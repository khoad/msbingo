// Package nbfs provides implementation of Microsoft [MC-NBFS]: .NET Binary Format: SOAP Data Structure
//
// More info https://msdn.microsoft.com/en-us/library/cc219175.aspx
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
