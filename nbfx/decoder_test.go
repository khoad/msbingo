package nbfx

import (
	//"io/ioutil"
	"testing"
	"bytes"
)

//https://golang.org/pkg/testing/

//func TestDecodeExample1(t *testing.T) {
//	decoder := NewDecoder()
//	path := "../examples/1"
//	bin, err := ioutil.ReadFile(path + ".bin")
//	if failOn(err, "unable to open "+path+".bin", t) {
//		return
//	}
//	_, err = decoder.Decode(bin)
//	if err == nil {
//		t.Error("Expected err")
//		return
//	} else if err.Error() != "Unknown Record ID 0x44" {
//		t.Error("Expected Unknown Record ID 0x44 message but got " + err.Error())
//		return
//	}
//}

func TestDecodePrefixDictionaryElementB(t *testing.T) {
	bin := []byte {0x45, 0x02}

	decoder := NewDecoder()
	actual, err := decoder.Decode(bin)
	if err != nil {
		t.Error("Unexpected error: " + err.Error() + " Got: " + actual)
		return
	}
	assertStringEqual(t, actual, "<b:str2>")
}

func TestPrefixDictionaryElementBName(t *testing.T) {
	codec := &codec{}
	record := getRecord(codec, 0x45)
	if record == nil {
		t.Error("Expected record but received nil")
		return
	} else if record.getName() != "PrefixDictionaryElementB (0x1)" {
		t.Error("Expected PrefixDictionaryElementB (0x1) but got " + record.getName())
	}
}

func TestDecodePrefixDictionaryElementAZ(t *testing.T) {
	bin := []byte{0x56, 0x02}

	decoder := NewDecoder()
	actual, err := decoder.Decode(bin)
	if err != nil {
		t.Error("Unexpected error: " + err.Error() + " Got: " + actual)
		return
	}
	assertStringEqual(t, actual, "<s:str2>")
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
	testReadMultiByteInt31(t, []byte{0x80, 0x80, 0x80}, 16384)
}

func TestReadMultiByteInt31_268435456(t *testing.T) {
	testReadMultiByteInt31(t, []byte{0x80, 0x80, 0x80, 0x01}, 268435456)
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
