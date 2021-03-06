package nbfx

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"encoding/binary"

	"github.com/satori/go.uuid"
)

type encoder struct {
	dict        map[string]uint32
	xml         *xml.Decoder
	bin         *bytes.Buffer
	tokenBuffer queue
}

func (e *encoder) addDictionaryString(index uint32, value string) {
	if _, ok := e.dict[value]; ok {
		return
	}
	e.dict[value] = index
}

// NewEncoder creates a new NBFX Encoder
func NewEncoder() Encoder {
	return NewEncoderWithStrings(nil)
}

// NewEncoderWithStrings creates a new NBFX Encoder with a dictionary (like an NBFS dictionary)
func NewEncoderWithStrings(dictionaryStrings map[uint32]string) Encoder {
	encoder := &encoder{dict: map[string]uint32{}}
	if dictionaryStrings != nil {
		for k, v := range dictionaryStrings {
			encoder.addDictionaryString(k, v)
		}
	}
	return encoder
}

func (e *encoder) popToken() (xml.Token, error) {
	var token xml.Token
	var err error
	if e.tokenBuffer.length > 0 {
		token = e.tokenBuffer.dequeue().(xml.Token)
		return token, err
	}
	token, err = e.xml.RawToken()
	token = xml.CopyToken(token) // make the token immutable (see doc for xml.Decoder.Token())
	return token, err
}

func (e *encoder) pushToken(token xml.Token) {
	e.tokenBuffer.enqueue(token)
}

func (e *encoder) Encode(reader io.Reader) ([]byte, error) {
	e.bin = &bytes.Buffer{}
	e.xml = xml.NewDecoder(reader)
	token, err := e.popToken()
	for err == nil && token != nil {
		record, err := e.getRecordFromToken(token)
		if err != nil {
			return e.bin.Bytes(), err
		}
		if record.isStartElement() {
			elementWriter := record.(elementRecordEncoder)
			err = elementWriter.encodeElement(e, token.(xml.StartElement))
		} else if record.isText() {
			textWriter := record.(textRecordEncoder)
			if _, ok := token.(xml.Comment); ok {
				err = textWriter.encodeText(e, textWriter, string(token.(xml.Comment)))
			} else {
				err = textWriter.encodeText(e, textWriter, string(token.(xml.CharData)))
			}
		} else if record.isEndElement() {
			elementWriter := record.(elementRecordEncoder)
			err = elementWriter.encodeElement(e, xml.StartElement{})
		} else {
			err = errors.New(fmt.Sprint("NotSupported: Encoding record", record))
		}
		if err != nil {
			return e.bin.Bytes(), fmt.Errorf("Error writing Token %s :: %s", token, err.Error())
		}
		token, err = e.popToken()
	}
	return e.bin.Bytes(), nil
}

func (e *encoder) getRecordFromToken(token xml.Token) (record, error) {
	switch token.(type) {
	case xml.StartElement:
		return e.getStartElementRecordFromToken(token.(xml.StartElement))
	case xml.CharData:
		return e.getTextRecordFromToken(token.(xml.CharData))
	case xml.EndElement:
		return records[endElement], nil
	case xml.Comment:
		return records[comment], nil
	}

	tokenXmlBytes, err := xml.Marshal(token)
	var tokenXml string
	if err != nil {
		tokenXml = "[[UNKNOWN]]"
	} else {
		tokenXml = string(tokenXmlBytes)
	}
	return nil, errors.New(fmt.Sprint("Unknown token", tokenXml))
}

func (e *encoder) getTextRecordFromToken(cd xml.CharData) (record, error) {
	withEndElement := false
	next, _ := e.popToken()
	switch next.(type) {
	case xml.EndElement:
		withEndElement = true
	}
	text := string(cd)
	if !withEndElement {
		e.pushToken(next)
	}
	return e.getTextRecordFromText(text, withEndElement)
}

func (e *encoder) getTextRecordFromText(text string, withEndElement bool) (record, error) {
	var id byte
	id = 0x00
	if text == "" {
		id = emptyText
	} else if text == "0" {
		id = zeroText
	} else if text == "1" {
		id = oneText
	} else if text == "false" {
		id = falseText
	} else if text == "true" {
		id = trueText
	} else if strings.Contains(text, " ") {
		id = startListText
	} else if isUuid(text) {
		id = uuidText
	} else if isUniqueId(text) {
		id = uniqueIdText
	} else if i, err := strconv.ParseInt(text, 10, 0); err == nil {
		if math.MinInt8 <= i && i <= math.MaxInt8 {
			id = int8Text
		} else if math.MinInt16 <= i && i <= math.MaxInt16 {
			id = int16Text
		} else if math.MinInt32 <= i && i <= math.MaxInt32 {
			id = int32Text
		} else if math.MinInt64 <= i && i <= math.MaxInt64 {
			id = int64Text
		} else {
			return nil, fmt.Errorf("Unknown integer record %v", i)
		}
	} else if _, err := strconv.ParseUint(text, 10, 0); err == nil {
		id = uInt64Text
	} else if _, err := strconv.ParseFloat(text, 64); err == nil {
		if isFloat32(text) {
			id = floatText
		} else {
			id = doubleText
		}
	} else if bSlice, err := b64.DecodeString(text); err == nil {
		lenBytes := len(bSlice)
		if lenBytes <= math.MaxUint8 {
			id = bytes8Text
		} else if lenBytes < math.MaxUint16 {
			id = bytes16Text
		} else if lenBytes < math.MaxUint32 {
			id = bytes32Text
		} else {
			return nil, fmt.Errorf("Base64 text too long, didn't encode: %v", text)
		}
	} else {
		if _, ok := e.dict[text]; ok || hasSpecialDictionaryPrefix(text) {
			id = dictionaryText
		} else if isQNameDictionaryText(text) {
			id = qNameDictionaryText
		} else {
			lenText := len(text)
			if lenText <= math.MaxUint8 {
				id = chars8Text
			} else if lenText < math.MaxUint16 {
				id = chars16Text
			} else if lenText < math.MaxUint32 {
				id = chars32Text
			} else {
				return nil, fmt.Errorf("Text too long, didn't encode: %v", text)
			}
		}
	}
	if id != 0 && withEndElement {
		id += 1
	}
	if rec, ok := records[id]; ok {
		return rec, nil
	}
	return nil, fmt.Errorf("Unknown text record id %#X for %s withEndElement %v", id, text, withEndElement)
}

func isQNameDictionaryText(text string) bool {
	if text[1] != ':' {
		return false
	}
	prefix := text[0]
	return 'a' <= prefix && prefix <= 'z'
}

func isFloat32(s string) bool {
	f3264, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return false
	}
	f32 := float32(f3264)
	return fmt.Sprint(f32) == s
}

func (e *encoder) getStartElementRecordFromToken(startElement xml.StartElement) (record, error) {
	prefix := startElement.Name.Space
	name := startElement.Name.Local
	prefixIndex := -1
	if len(prefix) == 1 && byte(prefix[0]) >= byte('a') && byte(prefix[0]) <= byte('z') {
		prefixIndex = int(byte(prefix[0]) - byte('a'))
	}
	isNameIndexAssigned := false
	if _, ok := e.dict[name]; ok {
		isNameIndexAssigned = true
	}
	localHasStrPrefix := hasSpecialDictionaryPrefix(startElement.Name.Local)

	if prefix == "" {
		if isNameIndexAssigned || localHasStrPrefix {
			return records[shortDictionaryElement], nil
		} else {
			return records[shortElement], nil
		}
	} else if prefixIndex != -1 {
		if isNameIndexAssigned || localHasStrPrefix {
			return records[prefixDictionaryElementA+byte(prefixIndex)], nil
		} else {
			return records[prefixElementA+byte(prefixIndex)], nil
		}
	} else {
		if isNameIndexAssigned || localHasStrPrefix {
			return records[dictionaryElement], nil
		} else {
			return records[element], nil
		}
	}
}

func (e *encoder) getAttributeRecordFromToken(attr xml.Attr) (record, error) {
	prefix := attr.Name.Space
	name := attr.Name.Local
	isXmlns := prefix == "xmlns" || prefix == "" && name == "xmlns"
	prefixIndex := -1
	if len(prefix) == 1 && byte(prefix[0]) >= byte('a') && byte(prefix[0]) <= byte('z') {
		prefixIndex = int(byte(prefix[0]) - byte('a'))
	}
	isNameIndexAssigned := false
	if _, ok := e.dict[name]; ok {
		isNameIndexAssigned = true
	}
	localHasStrPrefix := hasSpecialDictionaryPrefix(attr.Name.Local)
	valueHasStrPrefix := hasSpecialDictionaryPrefix(attr.Value)

	if prefix == "" {
		if isXmlns {
			if _, ok := e.dict[attr.Value]; ok || valueHasStrPrefix {
				return records[shortDictionaryXmlnsAttribute], nil
			} else {
				return records[shortXmlnsAttribute], nil
			}
		} else if isNameIndexAssigned || localHasStrPrefix {
			return records[shortDictionaryAttribute], nil
		} else {
			return records[shortAttribute], nil
		}
	} else if prefixIndex != -1 {
		if isNameIndexAssigned || localHasStrPrefix {
			return records[prefixDictionaryAttributeA+byte(prefixIndex)], nil
		} else {
			return records[prefixAttributeA+byte(prefixIndex)], nil
		}
	} else {
		if isXmlns {
			if isNameIndexAssigned || valueHasStrPrefix {
				return records[dictionaryXmlnsAttribute], nil
			} else {
				return records[xmlnsAttribute], nil
			}
		} else {
			if isNameIndexAssigned || localHasStrPrefix {
				return records[dictionaryAttribute], nil
			} else {
				return records[attribute], nil
			}
		}
	}
}

func hasSpecialDictionaryPrefix(str string) bool {
	return strings.HasPrefix(str, "str")
}

func writeString(e *encoder, str string) (int, error) {
	var strBytes = []byte(str)
	lenByteLen, err := writeMultiByteInt31(e, uint32(len(strBytes)))
	if err != nil {
		return lenByteLen, err
	}
	strByteLen, err := e.bin.Write(strBytes)
	return lenByteLen + strByteLen, err
}

func writeMultiByteInt31(e *encoder, num uint32) (int, error) {
	max := uint32(2147483647)
	if num > max {
		return 0, fmt.Errorf("Overflow: i (%d) must be <= max (%d)", num, max)
	}
	if num < maskMbi31 {
		return 1, e.bin.WriteByte(byte(num))
	}
	q := num / maskMbi31
	rem := num % maskMbi31
	err := e.bin.WriteByte(byte(maskMbi31 + rem))
	if err != nil {
		return 1, err
	}
	n, err := writeMultiByteInt31(e, q)
	return n + 1, err
}

func writeChars8Text(e *encoder, text string) error {
	bytes := []byte(text)
	writeMultiByteInt31(e, uint32(len(bytes)))
	_, err := e.bin.Write(bytes)
	return err
}

func writeChars16Text(e *encoder, text string) error {
	bytes := []byte(text)
	err := binary.Write(e.bin, binary.LittleEndian, uint16(len(bytes)))
	if err != nil {
		return err
	}
	_, err = e.bin.Write(bytes)
	return err
}

func writeChars32Text(e *encoder, text string) error {
	bytes := []byte(text)
	// PER SPEC: int32 NOT uint32
	err := binary.Write(e.bin, binary.LittleEndian, int32(len(bytes)))
	if err != nil {
		return err
	}
	_, err = e.bin.Write(bytes)
	return err
}

func writeUuidText(e *encoder, text string) error {
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

func writeDictionaryString(e *encoder, str string) error {
	if val, ok := e.dict[str]; ok {
		_, err := writeMultiByteInt31(e, val)
		if err != nil {
			return err
		}
	} else if hasSpecialDictionaryPrefix(str) {
		// capture "8" in "str8" and write "8"
		numString := str[3:]
		numInt, err := strconv.Atoi(numString)
		if err != nil {
			return err
		}
		_, err = writeMultiByteInt31(e, uint32(numInt))
		if err != nil {
			return err
		}
	} else {
		return errors.New("Invalid Operation")
	}
	return nil
}
