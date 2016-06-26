package nbfx

import (
	"fmt"
	"errors"
	"bytes"
)

type decoder struct {
	codec codec
}

func NewDecoder() Decoder {
	return NewDecoderWithStrings(nil)
}

func NewDecoderWithStrings(dictionaryStrings map[uint32]string) Decoder {
	decoder := &decoder{codec{make(map[uint32]string)}}
	if dictionaryStrings != nil {
		for k, v := range dictionaryStrings {
			decoder.codec.addDictionaryString(k, v)
		}
	}
	return decoder
}

func (d *decoder) Decode(bin []byte) (string, error) {
	reader := bytes.NewReader(bin)
	xml := bytes.Buffer{}
	b, err := reader.ReadByte()
	//println("ReadByte", string(b), err == nil)
	for err == nil {
		record := getRecord(&d.codec, b)
		//println("getRecord ", record)
		if record == nil {
			return "", errors.New(fmt.Sprintf("Unknown Record ID %#X", b))
		}
		bytes, err := record.read(reader)
		if err != nil {
			return "", err
		}
		xml.Write(bytes)
		b, err = reader.ReadByte()
	}
	return xml.String(), nil
}

type record interface {
	read(reader *bytes.Reader) ([]byte, error)
}

func getRecord(codec *codec, b byte) record {
	if b == 0x56 {
		return &prefixDictionaryElementS{codec}
	}
	return nil
}

type prefixDictionaryElementS struct {
	codec *codec
}

func (r *prefixDictionaryElementS) read(reader *bytes.Reader) ([]byte, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	key := uint32(b)
	if val, ok := r.codec.dict[key]; ok {
		return []byte(val), nil
	}
	return nil, errors.New(fmt.Sprint("Invalid DictionaryString str", key))
}
