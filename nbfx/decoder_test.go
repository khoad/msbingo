package nbfx

import (
	"bytes"
	//"io/ioutil"
	"testing"
)

//https://golang.org/pkg/testing/

func TestDecodeExampleEndElement(t *testing.T) {
	testDecode(t, []byte{0x40, 0x03, 0x64, 0x6F, 0x63, 0x01}, "<doc></doc>")
}

func TestDecodeExampleComment(t *testing.T) {
	testDecode(t, []byte{0x02, 0x07, 0x63, 0x6F, 0x6D, 0x6D, 0x65, 0x6E, 0x74}, "<!--comment-->")
}

func TestDecodePrefixDictionaryElementB(t *testing.T) {
	testDecode(t, []byte{0x45, 0x02}, "<b:str2>")
}

func TestDecodeOneText(t *testing.T) {
	testDecode(t, []byte{0x82}, "1")
}

func TestDecodePrefixDictionaryElementAZ(t *testing.T) {
	testDecode(t, []byte{0x56, 0x02}, "<s:str2>")
}

func TestReadMultiByteInt31_17(t *testing.T) {
	testReadMultiByteInt31(t, []byte{0x11}, 17)
}

func TestReadMultiByteInt31_145(t *testing.T) {
	testReadMultiByteInt31(t, []byte{0x91, 0x01}, 145)
}

func TestReadMultiByteInt31_5521(t *testing.T) {
	testReadMultiByteInt31(t, []byte{0x91, 0x2B}, 5521)
}

func TestReadMultiByteInt31_16384(t *testing.T) {
	testReadMultiByteInt31(t, []byte{0x80, 0x80, 0x01}, 16384)
}

func TestReadMultiByteInt31_2097152(t *testing.T) {
	testReadMultiByteInt31(t, []byte{0x80, 0x80, 0x80, 0x01}, 2097152)
}

func TestReadMultiByteInt31_268435456(t *testing.T) {
	testReadMultiByteInt31(t, []byte{0x80, 0x80, 0x80, 0x80, 0x01}, 268435456)
}

func TestReadString_abc(t *testing.T) {
	reader := bytes.NewReader([]byte{0x03, 0x61, 0x62, 0x63})
	actual, err := readString(reader)
	if err != nil {
		t.Error("Error: " + err.Error())
		return
	}
	expected := "abc"
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func testDecode(t *testing.T, bin []byte, expected string) {
	decoder := NewDecoder()
	actual, err := decoder.Decode(bin)
	if err != nil {
		t.Error("Unexpected error: " + err.Error() + " Got: " + actual)
	}
	assertStringEqual(t, actual, expected)
}

func testReadMultiByteInt31(t *testing.T, bin []byte, expected uint32) {
	reader := bytes.NewReader(bin)
	actual, err := readMultiByteInt31(reader)
	if err != nil {
		t.Error("Error: " + err.Error())
		return
	}
	if actual != expected {
		t.Errorf("Expected %d but got %d", expected, actual)
		return
	}
}

func assertStringEqual(t *testing.T, actual, expected string) {
	if expected != actual {
		t.Error(actual + " not equal to expected " + expected)
	}
}

func failOn(err error, message string, t *testing.T) bool {
	if err != nil {
		t.Error(message)
		return true
	}
	return false
}
