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
	return &decoder{}
}

func (d *decoder) Decode(bin []byte) (string, error) {
	reader := bytes.NewReader(bin)
	xml := bytes.Buffer{}
	b, err := reader.ReadByte()
	//println("ReadByte", string(b), err == nil)
	for err == nil {
		record := getRecord(b)
		//println("getRecord ", record)
		if record == nil {
			return "", errors.New(fmt.Sprintf("Unknown Record ID %x", b))
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

func getRecord(b byte) record {
	if b == 0x56 {
		return prefixDictionaryElementS{}
	}
	return nil
}

type prefixDictionaryElementS struct {
}

func (r prefixDictionaryElementS) read(reader *bytes.Reader) ([]byte, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	key := uint32(b)
	if val, ok := dict[key]; ok {
		return []byte(val), nil
	}
	return nil, errors.New(fmt.Sprint("Invalid DictionaryString str", key))
}
