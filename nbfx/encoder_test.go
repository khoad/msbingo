package nbfx

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"testing"
)

func TestEncodeExampleEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x01},
		"<doc></doc>")
}

func TestEncodeExampleComment(t *testing.T) {
	testEncode(t,
		[]byte{0x02, 0x07, 0x63, 0x6F, 0x6D, 0x6D, 0x65, 0x6E, 0x74},
		"<!--comment-->")
}

func TestEncodeExampleArray(t *testing.T) {
	testEncode(t,
		[]byte{0x03, 0x40, 0x03, 0x61, 0x72, 0x72, 0x01, 0x8B, 0x03, 0x33, 0x33, 0x88, 0x88, 0xDD, 0xDD},
		"<arr>13107</arr><arr>-30584</arr><arr>-8739</arr>")
}

func TestEncodeShortAttribute(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x04, 0x61, 0x74, 0x74, 0x72, 0x84, 0x01},
		"<doc attr=\"false\"></doc>")
}

func TestEncodeExampleAttribute(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x09, 0x03, 0x70, 0x72, 0x65, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x05, 0x03, 0x70, 0x72, 0x65, 0x04, 0x61, 0x74, 0x74, 0x72, 0x84, 0x01},
		"<doc xmlns:pre=\"http://abc\" pre:attr=\"false\"></doc>")
}

func TestEncodeExampleShortDictionaryAttribute(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0x08, 0x86, 0x01},
		"<doc str8=\"true\"></doc>")
}

func TestEncodeExampleDictionaryAttribute(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x09, 0x03, 0x70, 0x72, 0x65, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x07, 0x03, 0x70, 0x72, 0x65, 0x00, 0x86, 0x01},
		"<doc xmlns:pre=\"http://abc\" pre:str0=\"true\"></doc>")
}

func TestEncodeExampleShortXmlnsAttribute(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x08, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x01},
		"<doc xmlns=\"http://abc\"></doc>")
}

func TestEncodeExampleXmlnsAttribute(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x09, 0x01, 0x70, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x01},
		"<doc xmlns:p=\"http://abc\"></doc>")
}

func TestEncodeExampleShortDictionaryXmlnsAttribute(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x0A, 0x04, 0x01},
		"<doc xmlns=\"str4\"></doc>")
}

func TestEncodeExampleDictionaryXmlnsAttribute(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x0B, 0x01, 0x70, 0x04, 0x01},
		"<doc xmlns:p=\"str4\"></doc>")
}

func TestEncodeExamplePrefixDictionaryAttributeF(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x09, 0x01, 0x66, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x11, 0x0B, 0x98, 0x05, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x01},
		"<doc xmlns:f=\"http://abc\" f:str11=\"hello\"></doc>")
}

func TestEncodeExamplePrefixDictionaryAttributeX(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x09, 0x01, 0x78, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x23, 0x15, 0x98, 0x05, 0x77, 0x6F, 0x72, 0x6C, 0x64, 0x01},
		"<doc xmlns:x=\"http://abc\" x:str21=\"world\"></doc>")
}

func TestEncodeExamplePrefixAttributeK(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x09, 0x01, 0x6B, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x30, 0x04, 0x61, 0x74, 0x74, 0x72, 0x86, 0x01},
		"<doc xmlns:k=\"http://abc\" k:attr=\"true\"></doc>")
}

func TestEncodeExamplePrefixAttributeZ(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x09, 0x01, 0x7A, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x3F, 0x03, 0x61, 0x62, 0x63, 0x98, 0x03, 0x78, 0x79, 0x7A, 0x01},
		"<doc xmlns:z=\"http://abc\" z:abc=\"xyz\"></doc>")
}

func TestEncodeExampleShortElement(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x01},
		"<doc></doc>")
}

func TestEncodeExampleElement(t *testing.T) {
	testEncode(t,
		[]byte{0x41, 0x03, 0x70, 0x72, 0x65, 0x03, 0x64, 0x6F, 0x63, 0x09, 0x03, 0x70, 0x72, 0x65, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x01},
		"<pre:doc xmlns:pre=\"http://abc\"></pre:doc>")
}

func TestEncodeExampleShortDictionaryElement(t *testing.T) {
	testEncode(t,
		[]byte{0x42, 0x0E, 0x01},
		"<str14></str14>")
}

func TestEncodeExampleDictionaryElement(t *testing.T) {
	testEncode(t,
		[]byte{0x43, 0x03, 0x70, 0x72, 0x65, 0x0E, 0x09, 0x03, 0x70, 0x72, 0x65, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x01},
		"<pre:str14 xmlns:pre=\"http://abc\"></pre:str14>")
}

func TestEncodeExamplePrefixDictionaryElementA(t *testing.T) {
	testEncode(t,
		[]byte{0x44, 0x0A, 0x09, 0x01, 0x61, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x01},
		"<a:str10 xmlns:a=\"http://abc\"></a:str10>")
}

func TestEncodeExamplePrefixDictionaryElementS(t *testing.T) {
	testEncode(t,
		[]byte{0x56, 0x26, 0x09, 0x01, 0x73, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x01},
		"<s:str38 xmlns:s=\"http://abc\"></s:str38>")
}

func TestEncodeExamplePrefixElementA(t *testing.T) {
	testEncode(t,
		[]byte{0x5E, 0x05, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x09, 0x01, 0x61, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x01},
		"<a:hello xmlns:a=\"http://abc\"></a:hello>")
}

func TestEncodeExamplePrefixElementS(t *testing.T) {
	testEncode(t,
		[]byte{0x70, 0x09, 0x4D, 0x79, 0x4D, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x09, 0x01, 0x73, 0x0A, 0x68, 0x74, 0x74, 0x70, 0x3A, 0x2F, 0x2F, 0x61, 0x62, 0x63, 0x01},
		"<s:MyMessage xmlns:s=\"http://abc\"></s:MyMessage>")
}

func TestEncodeExampleZeroText(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0xA0, 0x03, 0x80, 0x01},
		"<doc str416=\"0\"></doc>")
}

func TestEncodeExampleZeroTextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x61, 0x62, 0x63, 0x81},
		"<abc>0</abc>")
}

func TestEncodeExampleOneText(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0x00, 0x82, 0x01},
		"<doc str0=\"1\"></doc>")
}

func TestEncodeExampleOneTextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x61, 0x62, 0x63, 0x83},
		"<abc>1</abc>")
}

func TestEncodeExampleFalseText(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0x00, 0x84, 0x01},
		"<doc str0=\"false\"></doc>")
}

func TestEncodeExampleFalseTextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x61, 0x62, 0x63, 0x85},
		"<abc>false</abc>")
}

func TestEncodeExampleTrueText(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0x00, 0x86, 0x01},
		"<doc str0=\"true\"></doc>")
}

func TestEncodeExampleTrueTextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x61, 0x62, 0x63, 0x87},
		"<abc>true</abc>")
}

func TestEncodeExampleInt8Text(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0xEC, 0x01, 0x88, 0xDE, 0x01},
		"<doc str236=\"-34\"></doc>")
}

func TestEncodeExampleInt8TextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x42, 0x9A, 0x01, 0x89, 0x7F},
		"<str154>127</str154>")
}

func TestEncodeExampleInt16Text(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0xEC, 0x01, 0x8A, 0x00, 0x80, 0x01},
		"<doc str236=\"-32768\"></doc>")
}

func TestEncodeExampleInt16TextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x42, 0x9A, 0x01, 0x8B, 0xFF, 0x7F},
		"<str154>32767</str154>")
}

func TestEncodeExampleInt32Text(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0xEC, 0x01, 0x8C, 0x15, 0xCD, 0x5B, 0x07, 0x01},
		"<doc str236=\"123456789\"></doc>")
}

func TestEncodeExampleInt32TextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x42, 0x9A, 0x01, 0x8D, 0xFF, 0xFF, 0xFF, 0x7F},
		"<str154>2147483647</str154>")
}

func TestEncodeExampleInt64Text(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0xEC, 0x01, 0x8E, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x01},
		"<doc str236=\"2147483648\"></doc>")
}

func TestEncodeExampleInt64TextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x42, 0x9A, 0x01, 0x8F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00},
		"<str154>1099511627776</str154>")
}

func TestEncodeExampleFloatText(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x01, 0x61, 0x90, 0xCD, 0xCC, 0x8C, 0x3F, 0x01},
		"<doc a=\"1.1\"></doc>")
}

func TestEncodeExampleFloatTextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x91, 0xCD, 0xCC, 0x01, 0x42},
		"<Price>32.45</Price>")
}

func TestEncodeExampleDoubleText(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x01, 0x61, 0x92, 0x74, 0x57, 0x14, 0x8B, 0x0A, 0xBF, 0x05, 0x40, 0x01},
		"<doc a=\"2.71828182845905\"></doc>")
}

func TestEncodeExampleDoubleTextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x02, 0x50, 0x49, 0x93, 0x11, 0x2D, 0x44, 0x54, 0xFB, 0x21, 0x09, 0x40},
		"<PI>3.14159265358979</PI>")
}

func TestEncodeExampleDecimalText(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x03, 0x69, 0x6E, 0x74, 0x94, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x2D, 0x4E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		"<doc int=\"5.123456\"></doc>")
}

func TestEncodeExampleDecimalTextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x08, 0x4D, 0x61, 0x78, 0x56, 0x61, 0x6C, 0x75, 0x65, 0x95, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		"<MaxValue>79228162514264337593543950335</MaxValue>")
}

func TestEncodeExampleDateTimeText(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0x6E, 0x96, 0xFF, 0x3F, 0x37, 0xF4, 0x75, 0x28, 0xCA, 0x2B, 0x01},
		"<doc str110=\"9999-12-31T23:59:59.9999999\"></doc>")
}

func TestEncodeExampleDateTimeTextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x42, 0x6C, 0x97, 0x00, 0x40, 0x8E, 0xF9, 0x5B, 0x47, 0xC8, 0x08},
		"<str108>2006-05-17T00:00:00</str108>")
}

func TestEncodeExampleChars8Text(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0x00, 0x98, 0x05, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x01},
		"<doc str0=\"hello\"></doc>")
}

func TestEncodeExampleChars8TextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x01, 0x61, 0x99, 0x05, 0x68, 0x65, 0x6C, 0x6C, 0x6F},
		"<a>hello</a>")
}

func TestEncodeExampleChars16Text(t *testing.T) {
	n := math.MaxUint8 + 2
	bytBuffer := bytes.NewBuffer([]byte{0x40, 0x01, 0x61, 0x06, 0x00, 0x9A})
	binary.Write(bytBuffer, binary.LittleEndian, uint16(n))
	strBuffer := bytes.Buffer{}
	strBuffer.WriteString("<a str0=\"")
	for i := 0; i < n; i++ {
		bytBuffer.WriteByte(0x62)
		strBuffer.WriteString("b")
	}
	bytBuffer.WriteByte(0x01)
	strBuffer.WriteString("\"></a>")
	testEncode(t,
		bytBuffer.Bytes(),
		strBuffer.String())
}

func TestEncodeExampleChars16TextWithEndElement(t *testing.T) {
	n := math.MaxUint8 + 2
	bytBuffer := bytes.NewBuffer([]byte{0x40, 0x01, 0x61, 0x9B})
	binary.Write(bytBuffer, binary.LittleEndian, uint16(n))
	strBuffer := bytes.Buffer{}
	strBuffer.WriteString("<a>")
	for i := 0; i < n; i++ {
		bytBuffer.WriteByte(0x62)
		strBuffer.WriteString("b")
	}
	strBuffer.WriteString("</a>")
	testEncode(t,
		bytBuffer.Bytes(),
		strBuffer.String())
}

func TestEncodeExampleChars32Text(t *testing.T) {
	n := math.MaxUint16 + 2
	bytBuffer := bytes.NewBuffer([]byte{0x40, 0x01, 0x61, 0x06, 0x00, 0x9C})
	binary.Write(bytBuffer, binary.LittleEndian, int32(n))
	strBuffer := bytes.Buffer{}
	strBuffer.WriteString("<a str0=\"")
	for i := 0; i < n; i++ {
		bytBuffer.WriteByte(0x62)
		strBuffer.WriteString("b")
	}
	bytBuffer.WriteByte(0x01)
	strBuffer.WriteString("\"></a>")
	testEncode(t,
		bytBuffer.Bytes(),
		strBuffer.String())
}

func TestEncodeExampleChars32TextWithEndElement(t *testing.T) {
	n := math.MaxUint16 + 2
	bytBuffer := bytes.NewBuffer([]byte{0x40, 0x01, 0x61, 0x9D})
	binary.Write(bytBuffer, binary.LittleEndian, int32(n))
	strBuffer := bytes.Buffer{}
	strBuffer.WriteString("<a>")
	for i := 0; i < n; i++ {
		bytBuffer.WriteByte(0x62)
		strBuffer.WriteString("b")
	}
	strBuffer.WriteString("</a>")
	testEncode(t,
		bytBuffer.Bytes(),
		strBuffer.String())
}

func TestEncodeExampleBytes8Text(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0x00, 0x9E, 0x08, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x01},
		"<doc str0=\"AAECAwQFBgc=\"></doc>")
}

func TestEncodeExampleBytes8TextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x06, 0x42, 0x61, 0x73, 0x65, 0x36, 0x34, 0x9F, 0x08, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
		"<Base64>AAECAwQFBgc=</Base64>")
}

func TestEncodeExampleBytes16(t *testing.T) {
	n := math.MaxUint8 + 3
	bytBuffer := bytes.NewBuffer([]byte{0x40, 0x01, 0x61, 0x06, 0x00, 0xA0})
	binary.Write(bytBuffer, binary.LittleEndian, uint16(n))
	strBuffer := bytes.Buffer{}
	strBuffer.WriteString("<a str0=\"")
	for i := 0; i < n; i++ {
		bytBuffer.WriteByte(0x05)
		if i%3 == 0 {
			strBuffer.WriteString("BQUF")
		}
	}
	bytBuffer.WriteByte(0x01)
	strBuffer.WriteString("\"></a>")
	testEncode(t,
		bytBuffer.Bytes(),
		strBuffer.String())
}

func TestEncodeExampleBytes16TextWithEndElement(t *testing.T) {
	n := math.MaxUint8 + 3
	bytBuffer := bytes.NewBuffer([]byte{0x40, 0x06, 0x42, 0x61, 0x73, 0x65, 0x36, 0x34, 0xA1})
	binary.Write(bytBuffer, binary.LittleEndian, uint16(n))
	strBuffer := bytes.Buffer{}
	strBuffer.WriteString("<Base64>")
	for i := 0; i < n; i++ {
		bytBuffer.WriteByte(0x05)
		if i%3 == 0 {
			strBuffer.WriteString("BQUF")
		}
	}
	strBuffer.WriteString("</Base64>")
	testEncode(t,
		bytBuffer.Bytes(),
		strBuffer.String())
}

func TestEncodeExampleBytes32(t *testing.T) {
	n := math.MaxUint16 + 3
	bytBuffer := bytes.NewBuffer([]byte{0x40, 0x01, 0x61, 0x06, 0x00, 0xA2})
	binary.Write(bytBuffer, binary.LittleEndian, int32(n))
	strBuffer := bytes.Buffer{}
	strBuffer.WriteString("<a str0=\"")
	for i := 0; i < n; i++ {
		bytBuffer.WriteByte(0x05)
		if i%3 == 0 {
			strBuffer.WriteString("BQUF")
		}
	}
	bytBuffer.WriteByte(0x01)
	strBuffer.WriteString("\"></a>")
	testEncode(t,
		bytBuffer.Bytes(),
		strBuffer.String())
}

func TestEncodeExampleBytes32TextWithEndElement(t *testing.T) {
	n := math.MaxUint16 + 3
	bytBuffer := bytes.NewBuffer([]byte{0x40, 0x06, 0x42, 0x61, 0x73, 0x65, 0x36, 0x34, 0xA3})
	binary.Write(bytBuffer, binary.LittleEndian, int32(n))
	strBuffer := bytes.Buffer{}
	strBuffer.WriteString("<Base64>")
	for i := 0; i < n; i++ {
		bytBuffer.WriteByte(0x05)
		if i%3 == 0 {
			strBuffer.WriteString("BQUF")
		}
	}
	strBuffer.WriteString("</Base64>")
	testEncode(t,
		bytBuffer.Bytes(),
		strBuffer.String())
}

func TestEncodeExampleStartListText(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x01, 0x61, 0xA4, 0x88, 0x7B, 0x98, 0x05, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x86, 0xA6, 0x01},
		"<doc a=\"123 hello true\"></doc>")
}

func TestEncodeExampleEndListText(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x01, 0x61, 0xA4, 0x88, 0x7B, 0x98, 0x05, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x86, 0xA6, 0x01},
		"<doc a=\"123 hello true\"></doc>")
}

func TestEncodeExampleEmptyText(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x01, 0x61, 0xA8, 0x01},
		"<doc a=\"\"></doc>")
}

func TestEncodeExampleEmptyTextWithEndElement(t *testing.T) {
	// This test is somewhat INVALID because when encoder sees an End Element,
	// it will encode the End Element (0x01) instead of the Empty Text with end element (0xA9)
	// In order to make this test pass like previously desired (encoding Empty Text with end
	// element 0xA9):

	//testEncode(t,
	//	[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0xA9},
	//	"<doc></doc>")

	// We would have to check for the previous Start Element and see if there is no attribute,
	// and check for the "Empty Text" which would be an End Element. Considering the whole
	// point of the codec is to save bytes, it doesn't save anymore bytes and introduce
	// unnecessary complexity. So for this test, let's have it encode a simple End Element
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x01},
		"<doc></doc>")
}

func TestEncodeExampleDictionaryText(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x02, 0x6E, 0x73, 0xAA, 0x38, 0x01},
		"<doc ns=\"str56\"></doc>")
}

func TestEncodeExampleDictionaryTextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x04, 0x54, 0x79, 0x70, 0x65, 0xAB, 0xC4, 0x01},
		"<Type>str196</Type>")
}

func TestEncodeExampleUniqueIdText(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0x00, 0xAC, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x01},
		"<doc str0=\"urn:uuid:33221100-5544-7766-8899-aabbccddeeff\"></doc>")
}

func TestEncodeExampleUniqueIdTextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x42, 0x1A, 0xAD, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF},
		"<str26>urn:uuid:33221100-5544-7766-8899-aabbccddeeff</str26>")
}

func TestEncodeExampleTimeSpanText(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0xAE, 0x00, 0xC4, 0xF5, 0x32, 0xFF, 0xFF, 0xFF, 0xFF, 0x01},
		"<doc>-PT5M44S</doc>")
}

func TestEncodeExampleTimeSpanTextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x42, 0x94, 0x07, 0xAF, 0x00, 0xB0, 0x8E, 0xF0, 0x1B, 0x00, 0x00, 0x00},
		"<str916>PT3H20M</str916>")
}

func TestEncodeExampleUuidText(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0x00, 0xB0, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x01},
		"<doc str0=\"03020100-0504-0706-0809-0a0b0c0d0e0f\"></doc>")
}

func TestEncodeExampleUuidTextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x02, 0x49, 0x44, 0xB1, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
		"<ID>03020100-0504-0706-0809-0a0b0c0d0e0f</ID>")
}

func TestEncodeExampleUInt64Text(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0x00, 0xB2, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x01},
		"<doc str0=\"18446744073709551615\"></doc>")
}

func TestEncodeExampleUInt64TextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x42, 0x9A, 0x01, 0xB3, 0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		"<str154>18446744073709551614</str154>")
}

func TestEncodeExampleBoolText(t *testing.T) {
	// Using true/false text is more efficient than bool text anyway

	//testEncode(t,
	//	[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0xB4, 0x01, 0x01},
	//	"<doc>true</doc>")

	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x87},
		"<doc>true</doc>")
}

func TestEncodeExampleBoolTextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x03, 0x40, 0x03, 0x61, 0x72, 0x72, 0x01, 0xB5, 0x05, 0x01, 0x00, 0x01, 0x00, 0x01},
		"<arr>true</arr><arr>false</arr><arr>true</arr><arr>false</arr><arr>true</arr>")
}

func TestEncodeExampleUnicodeChars8Text(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x01, 0x75, 0xB6, 0x06, 0x75, 0x00, 0x6E, 0x00, 0x69, 0x00, 0x01},
		"<doc u=\"uni\"></doc>")
}

func TestEncodeExampleUnicodeChars8TextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x01, 0x55, 0xB7, 0x06, 0x75, 0x00, 0x6E, 0x00, 0x69, 0x00},
		"<U>uni</U>")
}

func TestEncodeExampleUnicodeChars16Text(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x03, 0x75, 0x31, 0x36, 0xB8, 0x08, 0x00, 0x75, 0x00, 0x6E, 0x00, 0x69, 0x00, 0x32, 0x00, 0x01},
		"<doc u16=\"uni2\"></doc>")
}

func TestEncodeExampleUnicodeChars16TextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x55, 0x31, 0x36, 0xB9, 0x08, 0x00, 0x75, 0x00, 0x6E, 0x00, 0x69, 0x00, 0x32, 0x00},
		"<U16>uni2</U16>")
}

func TestEncodeExampleUnicodeChars32Text(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x04, 0x03, 0x75, 0x33, 0x32, 0xBA, 0x04, 0x00, 0x00, 0x00, 0x33, 0x00, 0x32, 0x00, 0x01},
		"<doc u32=\"32\"></doc>")
}

func TestEncodeExampleUnicodeChars32TextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x55, 0x33, 0x32, 0xBB, 0x04, 0x00, 0x00, 0x00, 0x33, 0x00, 0x32, 0x00},
		"<U32>32</U32>")
}

func TestEncodeExampleQNameDictionaryText(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x06, 0xF0, 0x06, 0xBC, 0x08, 0x8E, 0x07, 0x01},
		"<doc str880=\"i:str910\"></doc>")
}

func TestEncodeExampleQNameDictionaryTextWithEndElement(t *testing.T) {
	testEncode(t,
		[]byte{0x40, 0x04, 0x54, 0x79, 0x70, 0x65, 0xBD, 0x12, 0x90, 0x07},
		"<Type>s:str912</Type>")
}

//----------------------------------------------------

func TestEncodePrefixDictionaryElementB(t *testing.T) {
	xml := "<b:Foo>"

	encoder := NewEncoderWithStrings(map[uint32]string{0x02: "Foo"})
	actual, err := encoder.Encode(bytes.NewReader([]byte(xml)))
	if err != nil {
		t.Error("Unexpected error: " + err.Error() + " Got: " + fmt.Sprintf("%x", actual))
		return
	}
	assertBinEqual(t, actual, []byte{0x45, 0x02})
}

func TestEncodePrefixDictionaryElementS(t *testing.T) {
	xml := "<s:Foo>"

	encoder := NewEncoderWithStrings(map[uint32]string{0x02: "Foo"})
	actual, err := encoder.Encode(bytes.NewReader([]byte(xml)))
	if err != nil {
		t.Error("Unexpected error: " + err.Error() + " Got: " + fmt.Sprintf("%x", actual))
		return
	}
	assertBinEqual(t, actual, []byte{0x56, 0x02})
}

func TestWriteMultiByteInt31_17(t *testing.T) {
	testWriteMultiByteInt31(t, 17, []byte{0x11})
}

func TestWriteMultiByteInt31_145(t *testing.T) {
	testWriteMultiByteInt31(t, 145, []byte{0x91, 0x01})
}

func TestWriteMultiByteInt31_5521(t *testing.T) {
	testWriteMultiByteInt31(t, 5521, []byte{0x91, 0x2B})
}

func TestWriteMultiByteInt31_16384(t *testing.T) {
	testWriteMultiByteInt31(t, 16384, []byte{0x80, 0x80, 0x01})
}

func TestWriteMultiByteInt31_2097152(t *testing.T) {
	testWriteMultiByteInt31(t, 2097152, []byte{0x80, 0x80, 0x80, 0x01})
}

func TestWriteMultiByteInt31_268435456(t *testing.T) {
	testWriteMultiByteInt31(t, 268435456, []byte{0x80, 0x80, 0x80, 0x80, 0x01})
}

func TestWriteString_abc(t *testing.T) {
	buffer := bytes.Buffer{}
	str := "abc"
	expected := []byte{0x03, 0x61, 0x62, 0x63}
	expectedLen := len(expected)
	e := &encoder{bin: &buffer}
	i, err := writeString(e, str)
	if err != nil {
		t.Error("Error: " + err.Error())
		return
	}
	if i != expectedLen {
		t.Errorf("Expected to write %d bytes but wrote %d", expected, i)
	}
	assertBinEqual(t, buffer.Bytes(), expected)
}

func testWriteMultiByteInt31(t *testing.T, num uint32, expected []byte) {
	buffer := bytes.Buffer{}
	e := &encoder{bin: &buffer}
	i, err := writeMultiByteInt31(e, num)
	if err != nil {
		t.Error("Error: " + err.Error())
		return
	}
	if i != len(expected) {
		t.Errorf("Expected to write %d byte/s but wrote %d", len(expected), i)
		return
	}
	assertBinEqual(t, buffer.Bytes()[0:i], expected)
}

func testEncode(t *testing.T, expected []byte, xmlString string) {
	encoder := NewEncoder()
	actual, err := encoder.Encode(bytes.NewReader([]byte(xmlString)))
	if err != nil {
		t.Error("Unexpected error: " + err.Error() + " Got: " + string(actual))
	}
	assertBinEqual(t, actual, expected)
}

func assertBinEqual(t *testing.T, actual, expected []byte) {
	if len(actual) != len(expected) {
		t.Error("length of actual " + fmt.Sprint(len(actual)) + " not equal to length of expected " + fmt.Sprint(len(expected)))
	}
	for i, b := range actual {
		if i == len(expected) || b != expected[i] {
			pointerLine(t, actual, expected, i)
			return
		}
	}
	if len(actual) != len(expected) {
		pointerLine(t, actual, expected, len(actual))
	}
}

func pointerLine(t *testing.T, actual, expected []byte, i int) {
	pointerLine := "^^"
	for j := 1; j <= i; j++ {
		pointerLine = "--" + pointerLine
	}
	t.Error(fmt.Sprintf("actual\n%x\ndiffers from expected at index %d\n%x\n%s\n", actual, i, expected, pointerLine))
}
