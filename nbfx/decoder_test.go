package nbfx

import (
	"testing"
	"io/ioutil"
)

//https://golang.org/pkg/testing/

func TestDecodeExample1(t *testing.T) {
	decoder := NewDecoder()
	path := "../examples/1"
	bin, err := ioutil.ReadFile(path + ".bin")
	if failOn(err, "unable to open " + path + ".bin", t) { return }
	xmlBytes, err := ioutil.ReadFile(path + ".xml")
	if failOn(err, "unable to open " + path + ".xml", t) { return }
	expected := string(xmlBytes)
	actual := decoder.Decode(bin)
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
