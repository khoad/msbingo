package nbfx

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	. "math/big"
	"regexp"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

type decoder struct {
	dict         map[uint32]string
	elementStack stack
	bin          io.Reader
	xml          *xml.Encoder
}

func (d *decoder) addDictionaryString(index uint32, value string) {
	if _, ok := d.dict[index]; ok {
		return
	}
	d.dict[index] = value
}

// NewDecoder creates a new NBFX Decoder
func NewDecoder() Decoder {
	return NewDecoderWithStrings(nil)
}

// NewDecoderWithStrings creates a new NBFX Decoder with a dictionary (like an NBFS dictionary)
func NewDecoderWithStrings(dictionaryStrings map[uint32]string) Decoder {
	decoder := &decoder{dict: map[uint32]string{}}
	if dictionaryStrings != nil {
		for k, v := range dictionaryStrings {
			decoder.addDictionaryString(k, v)
		}
	}
	return decoder
}

func (d *decoder) Decode(reader io.Reader) (string, error) {
	// Use ioutil to read data from reader because if we try to read
	//  this manually, we can run into edge cases where the data we get
	//  back is unreliable and can contain extra unnecessary information
	//  as observed when reading through http response body containing
	//  extra zeros in sets of four.
	// This also seems to increase memory allocation efficiency
	// It is challenging to write a test for this bug as we haven't fully
	//  understood what the root cause for the extra zeros is.
	bytesRead, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	d.bin = bytes.NewBuffer(bytesRead)
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
	if uint32(b) < maskMbi31 {
		return uint32(b), err
	}
	nextB, err := readMultiByteInt31(reader)
	return maskMbi31*(nextB-1) + uint32(b), err
}

func readByte(reader io.Reader) (byte, error) {
	sb := make([]byte, 1)
	b, err := reader.Read(sb)
	if b > 0 {
		err = nil
	}
	return sb[0], err
}

func readStringBytes(reader io.Reader, len uint32) (string, error) {
	buf, err := readBytes(reader, len)
	return buf.String(), err
}

func readString(reader io.Reader) (string, error) {
	length, err := readMultiByteInt31(reader)
	if err != nil {
		return "", err
	}
	return readStringBytes(reader, length)
}

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
	var err error
	buf, err := readBytes(d.bin, 1)
	if err != nil {
		return "", err
	}
	var val uint8
	binary.Read(buf, binary.LittleEndian, &val)

	return readStringBytes(d.bin, uint32(val))
}

func readChars16Text(d *decoder) (string, error) {
	var err error
	buf, err := readBytes(d.bin, 2)
	if err != nil {
		return "", err
	}
	var val uint16
	binary.Read(buf, binary.LittleEndian, &val)

	return readStringBytes(d.bin, uint32(val))
}

func readChars32Text(d *decoder) (string, error) {
	var err error
	buf, err := readBytes(d.bin, 4)
	if err != nil {
		return "", err
	}
	var val uint32
	binary.Read(buf, binary.LittleEndian, &val)

	return readStringBytes(d.bin, val)
}

func readUnicodeChars8Text(d *decoder) (string, error) {
	var err error
	buf, err := readBytes(d.bin, 1)
	if err != nil {
		return "", err
	}
	var val uint8
	binary.Read(buf, binary.LittleEndian, &val)
	val /= 2

	return readUnicodeStringBytes(d.bin, uint32(val))
}

func readUnicodeChars16Text(d *decoder) (string, error) {
	var err error
	buf, err := readBytes(d.bin, 2)
	if err != nil {
		return "", err
	}
	var val uint16
	binary.Read(buf, binary.LittleEndian, &val)
	val /= 2

	return readUnicodeStringBytes(d.bin, uint32(val))
}

func readUnicodeChars32Text(d *decoder) (string, error) {
	var err error
	buf, err := readBytes(d.bin, 4)
	if err != nil {
		return "", err
	}
	var val uint32
	binary.Read(buf, binary.LittleEndian, &val)
	val /= 2

	return readUnicodeStringBytes(d.bin, val)
}

func readUnicodeStringBytes(r io.Reader, len uint32) (string, error) {
	runes := []rune{}
	for i := uint32(0); i < len; i++ {
		runeBuf, err := readBytes(r, 2)
		if err != nil {
			return string(runes), err
		}
		var runeInt int16
		binary.Read(runeBuf, binary.LittleEndian, &runeInt)
		theRune := rune(runeInt)
		runes = append(runes, theRune)
	}
	return string(runes), nil
}

func readBytes(reader io.Reader, numBytes uint32) (*bytes.Buffer, error) {
	var err error
	sb := make([]byte, numBytes)
	var b int
	b, err = reader.Read(sb)
	if b > 0 {
		err = nil
	}
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	_, err = buf.Write(sb)
	if err != nil {
		return nil, err
	}

	if uint32(b) < numBytes {
		nextBuf, err := readBytes(reader, numBytes-uint32(b))
		if err != nil {
			return &buf, err
		}
		_, err = buf.Write(nextBuf.Bytes())
		if err != nil {
			return &buf, err
		}
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
	d.bin.Read(make([]byte, 2)) // wReserved - ignored

	// scale - range 0 to 28
	buf, err := readBytes(d.bin, 1)
	if err != nil {
		return "", err
	}
	var scale byte
	binary.Read(buf, binary.LittleEndian, &scale)

	// sign: 0 = positive, 128 (0x80) = negative
	buf, err = readBytes(d.bin, 1)
	if err != nil {
		return "", err
	}
	var sign byte
	binary.Read(buf, binary.LittleEndian, &sign)

	// Hi32
	buf, err = readBytes(d.bin, 4)
	if err != nil {
		return "", err
	}
	var hi32 uint32
	binary.Read(buf, binary.LittleEndian, &hi32)

	// Lo64
	buf, err = readBytes(d.bin, 8)
	if err != nil {
		return "", err
	}
	var lo64 uint64
	binary.Read(buf, binary.LittleEndian, &lo64)

	// Goal: (Hi32 * 2^64 + Lo64) / 10^scale

	// 2^64
	var limit64 Int
	limit64.Exp(NewInt(2), NewInt(64), nil)

	var bigFirstPartInt Int
	// Hi32 * 2^64
	bigFirstPartInt.Mul(NewInt(int64(hi32)), &limit64)

	// Hi32 * 2^64 + Lo64
	bigFirstPartInt.Add(&bigFirstPartInt, new(Int).SetUint64(lo64))

	// (Hi32 * 2^64 + Lo64) / 10^scale
	numText := fmt.Sprint(&bigFirstPartInt)
	if scale > 0 {
		decIdx := len(numText) - int(scale)
		numText = numText[:decIdx] + "." + numText[decIdx:]
	}
	if sign == 0x80 {
		numText = "-" + numText
	}

	return numText, nil
}

func readDateTimeText(d *decoder) (string, error) {
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
	var nsec int64 = int64(cNanos % 1e9)

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
	buf, err := readBytes(d.bin, 8)
	if err != nil {
		return "", err
	}
	var val int64
	binary.Read(buf, binary.LittleEndian, &val)
	timeSpan := fmt.Sprint(time.Duration(val) * time.Nanosecond * 100)
	if strings.HasPrefix(timeSpan, "-") {
		timeSpan = strings.Replace(timeSpan, "-", "-PT", 1)
	} else {
		timeSpan = "PT" + timeSpan
	}
	timeSpan = strings.ToUpper(timeSpan)
	timeSpan = strings.Replace(timeSpan, "0S", "", 1)
	return timeSpan, nil
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
