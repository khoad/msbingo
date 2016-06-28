package nbfx

import (
	"io/ioutil"
	"testing"
)

//https://golang.org/pkg/testing/

func TestDecodeExample1(t *testing.T) {
	decoder := NewDecoder()
	path := "../examples/1"
	bin, err := ioutil.ReadFile(path + ".bin")
	if failOn(err, "unable to open "+path+".bin", t) {
		return
	}
	_, err = decoder.Decode(bin)
	if err == nil {
		t.Error("Expected err")
		return
	} else if err.Error() != "Unknown Record ID 0x44" {
		t.Error("Expected Unknown Record ID 0x44 message but got " + err.Error())
		return
	}
}

//func TestDecodePrefixDictionaryElementB(t *testing.T) {
//	bin := []byte {0x45, 0x02}
//
//	decoder := NewDecoder()
//	actual, err := decoder.Decode(bin)
//	if err != nil {
//		t.Error("Unexpected error: " + err.Error() + " Got: " + actual)
//		return
//	}
//	assertStringEqual(t, actual, "<b:str2>")
//}

func TestDecodePrefixDictionaryElementS(t *testing.T) {
	bin := []byte{0x56, 0x02}

	decoder := NewDecoder()
	actual, err := decoder.Decode(bin)
	if err != nil {
		t.Error("Unexpected error: " + err.Error() + " Got: " + actual)
		return
	}
	assertStringEqual(t, actual, "<s:str2>")
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
