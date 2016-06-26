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
		return t.Error("Expected err")
	} else if err.Error() != "Invalid DictionaryString value at offset 2" {
		return t.Error("Expected Invalid DictionaryString message")
	}
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

