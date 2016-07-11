package nbfx

import (
	"bytes"
	"encoding/xml"
	"io"
	"fmt"
	"encoding/binary"
	"math"
	"encoding/base64"
	"errors"
	"strings"
	//"time"
	//"github.com/nu7hatch/gouuid"
)

type decoder struct {
	dict        map[uint32]string
	elementStack Stack
	bin *bytes.Reader
	xml *xml.Encoder
}

func (d *decoder) addDictionaryString(index uint32, value string) {
	if _, ok := d.dict[index]; ok {
		return
	}
	d.dict[index] = value
}

func NewDecoder() Decoder {
	return NewDecoderWithStrings(nil)
}

func NewDecoderWithStrings(dictionaryStrings map[uint32]string) Decoder {
	decoder := &decoder{make(map[uint32]string), Stack{}, nil, nil}
	if dictionaryStrings != nil {
		for k, v := range dictionaryStrings {
			decoder.addDictionaryString(k, v)
		}
	}
	return decoder
}

func (d *decoder) Decode(bin []byte) (string, error) {
	d.bin = bytes.NewReader(bin)
	xmlBuf := &bytes.Buffer{}
	d.xml = xml.NewEncoder(xmlBuf)
	rec, err := getNextRecord(d)
	for err == nil && rec != nil {
		if rec.isStartElement() || rec.isEndElement() {
			//fmt.Println("Decoding", rec)
			elementReader := rec.(elementRecordDecoder)
			rec, err = elementReader.decodeElement(d)
		} else if rec.isText() {
			//fmt.Println("Decoding", rec)
			textReader := rec.(textRecordDecoder)
			_, err = textReader.decodeText(d, textReader)
			rec = nil
		} else {
			err = errors.New(fmt.Sprint("NotSupported: Decode record", rec))
		}
		if err == nil && rec == nil {
			rec, err = getNextRecord(d)
		}
	}
	//fmt.Println("Exiting main decoder loop...")
	d.xml.Flush()
	if err != nil && err != io.EOF {
		return xmlBuf.String(), err
	}
	return xmlBuf.String(), nil
}

func readMultiByteInt31(reader *bytes.Reader) (uint32, error) {
	b, err := reader.ReadByte()
	if uint32(b) < MASK_MBI31 {
		return uint32(b), err
	}
	nextB, err := readMultiByteInt31(reader)
	return MASK_MBI31*(nextB-1) + uint32(b), err
}

func readStringBytes(reader *bytes.Reader, readLenFunc func(r *bytes.Reader) (uint32, error)) (string, error) {
	len, err := readLenFunc(reader)
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

func readString(reader *bytes.Reader) (string, error) {
	return readStringBytes(reader, func(r *bytes.Reader) (uint32, error) {
		return readMultiByteInt31(r)
	})
}

var b64 = base64.StdEncoding.WithPadding(base64.StdPadding)

func readBytes8Text(d *decoder) (string, error) {
	var err error
	buf, err := readBytes(d.bin, 1)
	if err != nil {
		return "", err
	}
	var val uint8
	err = binary.Read(buf, binary.LittleEndian, &val)
	if err != nil {
		return "", err
	}
	buf, err = readBytes(d.bin, uint32(val))
	return b64.EncodeToString(buf.Bytes()), err
}

func readBytes16Text(d *decoder) (string, error) {
	var err error
	buf, err := readBytes(d.bin, 2)
	if err != nil {
		return "", err
	}
	var val uint16
	err = binary.Read(buf, binary.LittleEndian, &val)
	if err != nil {
		return "", err
	}
	buf, err = readBytes(d.bin, uint32(val))
	return b64.EncodeToString(buf.Bytes()), err
}

func readBytes32Text(d *decoder) (string, error) {
	var err error
	buf, err := readBytes(d.bin, 4)
	if err != nil {
		return "", err
	}
	var val uint32
	err = binary.Read(buf, binary.LittleEndian, &val)
	if err != nil {
		return "", err
	}
	buf, err = readBytes(d.bin, val)
	return b64.EncodeToString(buf.Bytes()), err
}

func readChars8Text(d *decoder) (string, error) {
	return readStringBytes(d.bin, func(r *bytes.Reader) (uint32, error) {
		var err error
		buf, err := readBytes(d.bin, 1)
		if err != nil {
			return uint32(0), err
		}
		var val uint8
		binary.Read(buf, binary.LittleEndian, &val)
		return uint32(val), err
	})
}

func readChars16Text(d *decoder) (string, error) {
	return readStringBytes(d.bin, func(r *bytes.Reader) (uint32, error) {
		var err error
		buf, err := readBytes(d.bin, 2)
		if err != nil {
			return uint32(0), err
		}
		var val uint16
		binary.Read(buf, binary.LittleEndian, &val)
		return uint32(val), err
	})
}

func readChars32Text(d *decoder) (string, error) {
	return readStringBytes(d.bin, func(r *bytes.Reader) (uint32, error) {
		var err error
		buf, err := readBytes(d.bin, 4)
		if err != nil {
			return uint32(0), err
		}
		var val uint32
		binary.Read(buf, binary.LittleEndian, &val)
		return uint32(val), err
	})
}

func readUnicodeChars8Text(d *decoder) (string, error) {
	return readUnicodeStringBytes(d.bin, func(r *bytes.Reader) (uint32, error) {
		var err error
		buf, err := readBytes(d.bin, 1)
		if err != nil {
			return uint32(0), err
		}
		var val uint8
		binary.Read(buf, binary.LittleEndian, &val)
		val /= 2
		return uint32(val), err
	})
}

func readUnicodeChars16Text(d *decoder) (string, error) {
	return readUnicodeStringBytes(d.bin, func(r *bytes.Reader) (uint32, error) {
		var err error
		buf, err := readBytes(d.bin, 2)
		if err != nil {
			return uint32(0), err
		}
		var val uint8
		binary.Read(buf, binary.LittleEndian, &val)
		val /= 2
		return uint32(val), err
	})
}

func readUnicodeChars32Text(d *decoder) (string, error) {
	return readUnicodeStringBytes(d.bin, func(r *bytes.Reader) (uint32, error) {
		var err error
		buf, err := readBytes(d.bin, 4)
		if err != nil {
			return uint32(0), err
		}
		var val uint8
		binary.Read(buf, binary.LittleEndian, &val)
		val /= 2
		return uint32(val), err
	})
}

func readUnicodeStringBytes(r *bytes.Reader, readLenFunc func(r *bytes.Reader) (uint32, error)) (string, error) {
	len, err := readLenFunc(r)
	if err != nil {
		return "", err
	}
	runes := []rune{}
	for i := uint32(0); i < len; {
		runeBuf, err := readBytes(r, 2)
		if err != nil {
			return string(runes), err
		}
		var runeInt int16
		binary.Read(runeBuf, binary.LittleEndian, &runeInt)
		theRune := rune(runeInt)
		runes = append(runes, theRune)
		i++
	}
	return string(runes), nil
}

func readBytes(reader *bytes.Reader, numBytes uint32) (*bytes.Buffer, error) {
	var err error
	buf := &bytes.Buffer{}
	var i uint32
	for i = 0; i < numBytes && err == nil; i++ {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		buf.WriteByte(b)
		if err != nil {
			return nil, err
		}
	}
	return buf, nil
}

func readInt8Text(d *decoder) (string, error) {
	var err error
	buf, err := readBytes(d.bin, 1)
	if err != nil {
		return "", err
	}
	var val int8
	binary.Read(buf, binary.LittleEndian, &val)
	return fmt.Sprintf("%d", val), nil
}

func readInt16Text(d *decoder) (string, error) {
	var err error
	buf, err := readBytes(d.bin, 2)
	if err != nil {
		return "", err
	}
	var val int16
	binary.Read(buf, binary.LittleEndian, &val)
	return fmt.Sprintf("%d", val), nil
}

func readInt32Text(d *decoder) (string, error) {
	var err error
	buf, err := readBytes(d.bin, 4)
	if err != nil {
		return "", err
	}
	var val int32
	binary.Read(buf, binary.LittleEndian, &val)
	return fmt.Sprintf("%d", val), nil
}

func readInt64Text(d *decoder) (string, error) {
	var err error
	buf, err := readBytes(d.bin, 8)
	if err != nil {
		return "", err
	}
	var val int64
	binary.Read(buf, binary.LittleEndian, &val)
	return fmt.Sprintf("%d", val), nil
}

func readUInt64Text(d *decoder) (string, error) {
	var err error
	buf, err := readBytes(d.bin, 8)
	if err != nil {
		return "", err
	}
	var val uint64
	binary.Read(buf, binary.LittleEndian, &val)
	return fmt.Sprintf("%d", val), nil
}

func readFloatText(d *decoder) (string, error) {
	var err error
	buf, err := readBytes(d.bin, 4)
	if err != nil {
		return "", err
	}
	var val float32
	binary.Read(buf, binary.LittleEndian, &val)
	if val == float32(math.Inf(1)) {
		return "INF", nil
	} else if val == float32(math.Inf(-1)) {
		return "-INF", nil
	}
	return fmt.Sprintf("%v", val), nil
}

func readDoubleText(d *decoder) (string, error) {
	var err error
	buf, err := readBytes(d.bin, 8)
	if err != nil {
		return "", err
	}
	var val float64
	binary.Read(buf, binary.LittleEndian, &val)
	if val == math.Inf(1) {
		return "INF", nil
	} else if val == math.Inf(-1) {
		return "-INF", nil
	}
	return fmt.Sprintf("%v", val), nil
}

func readListText(d *decoder) (string, error) {
	items := []string{}
	for {
		rec, err := getNextRecord(d)
		if err != nil {
			return "", err
		}
		if !rec.isText() {
			return "", errors.New("Records within list must be TextRecord types, but found " + rec.getName())
		}
		if rec.getName() == records[EndListText].getName() {
			break
		}
		textDecoder := rec.(textRecordDecoder)
		item, err := textDecoder.readText(d)
		if err != nil {
			return "", err
		}
		items = append(items, item)
	}
	return strings.Join(items, " "), nil
}

func readDecimalText(d *decoder) (string, error) {
	return "", errors.New("NotImplemented: DecimalText")
}

func readDateTimeText(d *decoder) (string, error) {
	//bin, err := readBytes(reader, 8)
	//if err != nil {
	//	return "", err
	//}
	//timeBin := bin[:]
	//tzBin := bin[7]
	//tzBin[7] = tzBin[7] & 0xFC
	//tzBin = tzBin & 0x03
	//timeUint := uint64(timeBin)
	//time := time.Time(timeUint)
	//time.
	return "", errors.New("NotImplemented: DateTimeText")
}

func readUniqueIdText(d *decoder) (string, error) {
	//var err error
	//buf, err := readBytes(reader, 16)
	//if err != nil {
	//	return "", err
	//}
	//val, err := uuid.Parse(buf.Bytes())
	//if err != nil {
	//	return "", err
	//}
	//val.setVariant(uuid.ReservedRFC4122)
	//return fmt.Sprintf("%s", val.String()), nil
	return "", errors.New("NotImplemented: UniqueIdText")
}

func readUuidText(d *decoder) (string, error) {
	return "", errors.New("NotImplemented: UuidText")
}

func readTimeSpanText(d *decoder) (string, error) {
	return "", errors.New("NotImplemented: TimeSpanText")
}

func readBoolText(d *decoder) (string, error) {
	b, err := d.bin.ReadByte()
	if err != nil {
		return "", err
	}
	if b == 0 {
		return "false", nil
	} else if b == 1 {
		return "true", nil
	}
	return "", errors.New("BoolText record byte must be 0 or 1")
}

func readQNameDictionaryText(d *decoder) (string, error) {
	b, err := d.bin.ReadByte()
	if err != nil {
		return "", err
	}
	prefix := string('a' + b)
	name, err := readDictionaryString(d)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", prefix, name), nil
}

func readDictionaryString(d *decoder) (string, error) {
	key, err := readMultiByteInt31(d.bin)
	if err != nil {
		return "", err
	}
	if val, ok := d.dict[key]; ok {
		return val, nil
	}
	return fmt.Sprintf("str%d", key), nil
}
