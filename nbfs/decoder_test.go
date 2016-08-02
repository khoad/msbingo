package nbfs

import (
	"bytes"
	"io/ioutil"
	"testing"
	"net/http"
)

func TestDecodeExample1(t *testing.T) {
	decoder := NewDecoder()
	path := "../examples/1"
	bin, err := ioutil.ReadFile(path + ".bin")
	if failOn(err, "unable to open "+path+".bin", t) {
		return
	}
	xmlBytes, err := ioutil.ReadFile(path + ".xml")
	if failOn(err, "unable to open "+path+".xml", t) {
		return
	}
	expected := string(xmlBytes)
	actual, err := decoder.Decode(bytes.NewReader(bin))
	if err != nil {
		t.Error("Unexpected error: " + err.Error() + " Got: " + actual)
		return
	}
	if expected != actual {
		t.Error("actual\n" + actual + "\nnot equal to expected\n" + expected)
	}
	assertEqual(t, actual, expected)
}

func TestDecodeExample1ThroughHttpServer(t *testing.T) {
	path := "../examples/1"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bin, err := ioutil.ReadFile(path + ".bin")
		if failOn(err, "unable to open "+path+".bin", t) {
			return
		}
		w.Write(bin)
	})
	go http.ListenAndServe(":12345", nil)

	resp, _ := http.Get("http://localhost:12345")
	decoder := NewDecoder()

	actual, err := decoder.Decode(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		t.Error("Unexpected error: " + err.Error() + " Got: " + actual)
		return
	}
	xmlBytes, err := ioutil.ReadFile(path + ".xml")
	if failOn(err, "unable to open "+path+".xml", t) {
		return
	}
	expected := string(xmlBytes)
	assertEqual(t, actual, expected)
}

func TestDecodePrefixDictionaryElementS(t *testing.T) {
	bin := []byte{0x56, 0x02}

	decoder := NewDecoder()
	actual, err := decoder.Decode(bytes.NewReader(bin))
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
