package nbfx

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"io"
	"math"
	"regexp"
	"strings"
	"time"
)

type decoder struct {
	dict         map[uint32]string
	elementStack Stack
	bin          io.Reader
	xml          *xml.Encoder
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

func (d *decoder) Decode(reader io.Reader) (string, error) {
	d.bin = reader
	xmlBuf := &bytes.Buffer{}
	d.xml = xml.NewEncoder(xmlBuf)
	rec, err := getNextRecord(d)
	for err == nil && rec != nil {
		if rec.isStartElement() || rec.isEndElement() {
			elementReader := rec.(elementRecordDecoder)
			rec, err = elementReader.decodeElement(d)
		} else if rec.isText() {
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
	d.xml.Flush()
	if err != nil && err != io.EOF {
		return xmlBuf.String(), err
	}
	return xmlBuf.String(), nil
}

func readMultiByteInt31(reader io.Reader) (uint32, error) {
	b, err := readByte(reader)
	if uint32(b) < mask_mbi31 {
		return uint32(b), err
	}
	nextB, err := readMultiByteInt31(reader)
	return mask_mbi31*(nextB-1) + uint32(b), err
}

func readByte(reader io.Reader) (byte, error) {
	sb := make([]byte, 1)
	_, err := reader.Read(sb)
	return sb[0], err
}

func readStringBytes(reader io.Reader, readLenFunc func(r io.Reader) (uint32, error)) (string, error) {
	len, err := readLenFunc(reader)
	if err != nil {
		return "", err
	}

	buf, err := readBytes(reader, len)
	return buf.String(), err
}

func readString(reader io.Reader) (string, error) {
	return readStringBytes(reader, func(r io.Reader) (uint32, error) {
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
	return readStringBytes(d.bin, func(r io.Reader) (uint32, error) {
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
	return readStringBytes(d.bin, func(r io.Reader) (uint32, error) {
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
	return readStringBytes(d.bin, func(r io.Reader) (uint32, error) {
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
	return readUnicodeStringBytes(d.bin, func(r io.Reader) (uint32, error) {
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
	return readUnicodeStringBytes(d.bin, func(r io.Reader) (uint32, error) {
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
	return readUnicodeStringBytes(d.bin, func(r io.Reader) (uint32, error) {
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

func readUnicodeStringBytes(r io.Reader, readLenFunc func(r io.Reader) (uint32, error)) (string, error) {
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

func readBytes(reader io.Reader, numBytes uint32) (*bytes.Buffer, error) {
	var err error
	sb := make([]byte, numBytes)
	_, err = reader.Read(sb)
	if err != nil && err != io.EOF {
		return nil, err
	}

	buf := bytes.Buffer{}
	_, err = buf.Write(sb)
	if err != nil {
		return nil, err
	}
	return &buf, nil
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
		if rec.getName() == records[endListText].getName() {
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
	d.bin.Read(make([]byte, 16))
	return "[DECIMAL]", nil
	//return "", errors.New("NotImplemented: DecimalText")
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

	// https://play.golang.org/p/Hy9NNuD7u5
	//d.bin.Read(make([]byte, 8))
	//return "[DATETIME]", nil
	//return "", errors.New("NotImplemented: DateTimeText")

	buf, err := readBytes(d.bin, 8)
	if err != nil {
		return "", err
	}

	bin := buf.Bytes()
	tz := (bin[7] & 0xC0) >> 6

	// Masking
	bin[7] &= 0x3F
	// Create a new buffer on the new masked bin
	buf = bytes.NewBuffer(bin)

	var maskedUIntDate uint64
	binary.Read(buf, binary.LittleEndian, &maskedUIntDate)

	// cNanos for cent-nanos (NBFX spec states the number is the 100 nanoseconds that have elapsed since 1.1.0001)
	var cNanos uint64 = maskedUIntDate
	var sec int64 = int64(cNanos / 1e7)
	var nsec int64 = int64(cNanos % 1e7)

	const (
		secondsPerDay        = 24 * 60 * 60
		unixToInternal int64 = (1969*365 + 1969/4 - 1969/100 + 1969/400) * secondsPerDay
		internalToUnix int64 = -unixToInternal
	)
	t := time.Unix(sec+internalToUnix, nsec)

	switch tz {
	case 0: // not specified, nothing is added
		return fmt.Sprint(t.UTC().Format("2006-01-02T15:04:05.9999999")), nil
	case 1: // UTC, add "Z"
		// TODO: This won't return what we want
		return fmt.Sprint(t.UTC()), nil
	case 2: // Local, add offset
		// TODO: This won't return what we want
		return fmt.Sprint(t.Local()), nil
	}

	return "", fmt.Errorf("Unrecognized TZ %v", tz)
}

func readUniqueIdText(d *decoder) (string, error) {
	result, err := readUuidText(d)
	if err != nil {
		return "", err
	}
	return urnPrefix + result, nil
}

const urnPrefix string = "urn:uuid:"

func isUniqueId(text string) bool {

	if !strings.HasPrefix(text, urnPrefix) {
		return false
	}
	uuidString := text[len(urnPrefix):]
	return isUuid(uuidString)
}

func flipBytes(bin []byte) []byte {
	for i, j := 0, len(bin)-1; i < j; i, j = i+1, j-1 {
		bin[i], bin[j] = bin[j], bin[i]
	}

	return bin
}

func flipUuidByteOrder(bin []byte) ([]byte, error) {
	part1 := flipBytes(bin[0:4])
	part2 := flipBytes(bin[4:6])
	part3 := flipBytes(bin[6:8])
	part4 := bin[8:]

	//concatenate parts 1-4
	return append(part1, append(part2, append(part3, part4...)...)...), nil
}

func writeUniqueIdText(e *encoder, text string) error {
	id, err := uuid.FromString(text)
	bin := id.Bytes()
	bin, err = flipUuidByteOrder(bin)
	if err != nil {
		return err
	}
	_, err = e.bin.Write(bin)
	if err != nil {
		return err
	}
	return nil
}

func readUuidText(d *decoder) (string, error) {
	bytes := make([]byte, 16)

	_, err := d.bin.Read(bytes)
	if err != nil {
		return "", err
	}

	bytes, err = flipUuidByteOrder(bytes)
	if err != nil {
		return "", err
	}

	val, err := uuid.FromBytes(bytes)
	if err != nil {
		return "", err
	}

	return val.String(), nil
}

func isUuid(text string) bool {
	if len(text) != 36 {
		return false
	}
	match, err := regexp.MatchString("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", text)
	if err != nil {
		return false
	}
	return match
}

func readTimeSpanText(d *decoder) (string, error) {
	return "", errors.New("NotImplemented: TimeSpanText")
}

func readBoolText(d *decoder) (string, error) {
	b, err := readByte(d.bin)
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
	b, err := readByte(d.bin)
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
