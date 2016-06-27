package nbfs

import (
	"testing"
	"io/ioutil"
)

func TestDecodeExample1(t *testing.T) {
	decoder := NewDecoder()
	path := "../examples/1"
	bin, err := ioutil.ReadFile(path + ".bin")
	if failOn(err, "unable to open " + path + ".bin", t) { return }
	xmlBytes, err := ioutil.ReadFile(path + ".xml")
	if failOn(err, "unable to open " + path + ".xml", t) { return }
	expected := string(xmlBytes)
	actual, err := decoder.Decode(bin)
	if err != nil {
		t.Error("Unexpected error: " + err.Error() + " Got: " + actual)
		return
	}
	if expected != actual {
		t.Error(actual + "\nnot equal to expected\n" + expected)
	}
	assertEqual(t, actual, expected)
}

func TestDecodePrefixDictionaryElementS(t *testing.T) {
	bin := []byte {0x56, 0x02}

	decoder := NewDecoder()
	actual, err := decoder.Decode(bin)
	if err != nil {
		t.Error("Unexpected error: " + err.Error() + " Got: " + actual)
		return
	}
	assertEqual(t, actual, "<s:Envelope>")
}

func assertEqual(t *testing.T, actual, expected string) {
	if expected != actual {
		t.Error(actual + "\nnot equal to expected\n" + expected)
	}
}

func failOn(err error, message string, t *testing.T) bool {
	if err != nil {
		t.Error(message)
		return true
	}
	return false
}

