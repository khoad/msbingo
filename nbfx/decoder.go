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
		//fmt.Println("Record:", record.name())
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
	name() string
}

func getRecord(codec *codec, b byte) record {
	if b == 0x56 {
		return &prefixDictionaryElementS{codec}
	} else if b == 0x0B {
		return &dictionaryXmlnsAttribute{codec}
	}
	return nil
}

//(0x56)
type prefixDictionaryElementS struct {
	codec *codec
}

func (r *prefixDictionaryElementS) name() string {
	return "PrefixDictionaryElementS (0x56)"
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

//(0x0B)
type dictionaryXmlnsAttribute struct {
	codec *codec
}

func (r *dictionaryXmlnsAttribute) name() string {
	return "dictionaryXmlnsAttribute (0x0B)"
}

func (r *dictionaryXmlnsAttribute) read(reader *bytes.Reader) ([]byte, error) {
	strBytes, err := readString(reader)
	if err != nil {
		return strBytes, err
	}
	strBytes = append(strBytes, []byte(":")[0])
	_, err = reader.ReadByte()
	if err != nil {
		return strBytes, err
	}
	return strBytes, nil
}

func readMultiByteInt31(reader * bytes.Reader) (uint32, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	return uint32(b), nil //TODO: Handle multibyte values!!!
}

func readString(reader *bytes.Reader) ([]byte, error) {
	var len uint32
	len, err := readMultiByteInt31(reader)
	if err != nil {
		return nil, err
	}
	strBytes := *new([]byte)
	for i := uint32(0); i < len; {
		b, err := reader.ReadByte()
		if err != nil {
			return strBytes, err
		}
		strBytes = append(strBytes, b)
		i++
	}
	return strBytes, nil
}
