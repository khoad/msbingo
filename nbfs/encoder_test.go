package nbfs

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestEncodeExample1(t *testing.T) {
	encoder := NewEncoder()
	path := "../examples/1"
	xmlBin, err := ioutil.ReadFile(path + ".xml")
	if failOn(err, "unable to open "+path+".xml", t) {
		return
	}
	expected, err := ioutil.ReadFile(path + ".bin")
	if failOn(err, "unable to open "+path+".bin", t) {
		return
	}
	actual, err := encoder.Encode(string(xmlBin))
	if err != nil {
		t.Error(fmt.Sprint(err.Error() + " : Got ", actual))
		return
	}
	assertBinEqual(t, actual, expected)
}

func assertBinEqual(t *testing.T, actual, expected []byte) {
	if len(actual) != len(expected) {
		t.Error("length of actual " + fmt.Sprint(len(actual)) + " not equal to length of expected " + fmt.Sprint(len(expected)))
	}
	for i, b := range actual {
		if b != expected[i] {
			t.Error(fmt.Sprintf("actual\n%x\ndiffers from expected\n%x\nat index %d", actual, expected, i))
			return
		}
	}
}
