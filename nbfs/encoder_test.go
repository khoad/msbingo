package nbfs

import (
	"bytes"
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
	actual, err := encoder.Encode(bytes.NewReader(xmlBin))
	if err != nil {
		t.Error(fmt.Sprint(err.Error()+" : Got ", actual))
		return
	}
	assertBinEqual(t, actual, expected)
}
