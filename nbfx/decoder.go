package nbfx

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
)

type decoder struct {
	codec codec
}

func NewDecoder() Decoder {
	return NewDecoderWithStrings(nil)
}

func NewDecoderWithStrings(dictionaryStrings map[uint32]string) Decoder {
	decoder := &decoder{codec{make(map[uint32]string), make(map[string]uint32)}}
	if dictionaryStrings != nil {
		for k, v := range dictionaryStrings {
			decoder.codec.addDictionaryString(k, v)
		}
	}
	return decoder
}

func (d *decoder) Decode(bin []byte) (string, error) {
	reader := bytes.NewReader(bin)
	xmlBuf := &bytes.Buffer{}
	xmlEncoder := xml.NewEncoder(xmlBuf)
	b, err := reader.ReadByte()
	var startingElement xml.StartElement
	haveStartingElement := false
	flushStartElement := func() {
		if haveStartingElement {
			xmlEncoder.EncodeToken(startingElement)
		}
		haveStartingElement = false
		startingElement = xml.StartElement{}
	}
	initStartElement := func(token xml.Token) {
		flushStartElement()
		haveStartingElement = true
		startingElement = token.(xml.StartElement)
	}
	for err == nil {
		record := getRecord(&d.codec, b)
		if record == nil {
			xmlEncoder.Flush()
			return xmlBuf.String(), errors.New(fmt.Sprintf("Unknown Record ID %#x", b))
		}
		var token xml.Token
		token, err = record.read(reader)
		if err != nil {
			xmlEncoder.Flush()
			return xmlBuf.String(), err
		}
		if record.isElementStart() {
			initStartElement(token)
		} else if record.isAttribute() {
			startingElement.Attr = append(startingElement.Attr, token.(xml.Attr))
		} else {
			flushStartElement()
			xmlEncoder.EncodeToken(token)
		}

		b, err = reader.ReadByte()
	}
	flushStartElement()
	xmlEncoder.Flush()
	if err != nil && err != io.EOF {
		return xmlBuf.String(), err
	}
	return xmlBuf.String(), nil
}

var MASK_MBI31 = byte(0x80)

func readMultiByteInt31(reader *bytes.Reader) (uint32, error) {
	buf := new([5]byte)
	keepReading := true
	i := -1
	for keepReading {
		i++
		b, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}
		if b >= MASK_MBI31 {
			b -= MASK_MBI31
			keepReading = true
		} else {
			keepReading = false
		}
		buf[i] = b
	}
	var val uint32
	for ; i >= 0; i-- {
		val += uint32(buf[i]) * uint32(math.Pow(128, float64(i)))
	}
	return val, nil
}

func writeMultiByteInt31(writer io.Writer, num uint32) (int, error) {
	max := uint32(2147483647)
	if num > max {
		return 0, errors.New(fmt.Sprintf("Overflow: i (%d) must be <= max (%d)", num, max))
	}
	buf := new([5]byte)
	val := num
	i := 4
	lastByte := 0
	for ; i >= 0; i-- {
		var base uint32
		if i > 0 {
			base = uint32(math.Pow(128, float64(i)))
		} else {
			base = 0
		}
		digit := byte(0x00)
		if val >= base {
			if base > 0 {
				digit = byte(math.Floor(float64(val / base)))
				val -= uint32(digit) * base
			} else {
				digit = byte(val)
			}
		}
		buf[i] = digit
	}

	haveLastByte := false
	for j := len(buf) - 1; j >= 0; j-- {
		if !haveLastByte && buf[j] > 0x00 {
			haveLastByte = true
			lastByte = j
		} else if haveLastByte {
			buf[j] = buf[j] + MASK_MBI31
		}
	}

	return writer.Write(buf[0 : lastByte+1])
}

func readString(reader *bytes.Reader) (string, error) {
	var len uint32
	len, err := readMultiByteInt31(reader)
	if err != nil {
		return "", err
	}
	strBuffer := bytes.Buffer{}
	for i := uint32(0); i < len; {
		b, err := reader.ReadByte()
		if err != nil {
			return strBuffer.String(), err
		}
		strBuffer.WriteByte(b)
		i++
	}
	return strBuffer.String(), nil
}
