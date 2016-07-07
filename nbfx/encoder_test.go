package nbfx

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
	actual, err := encoder.Encode(string(xmlBin))
	if err != nil {
		t.Error(err.Error())
		return
	}
	assertBinEqual(t, actual, expected)
}

func TestEncodePrefixDictionaryElementB(t *testing.T) {
	xml := "<b:Foo>"

	encoder := NewEncoderWithStrings(map[uint32]string{0x02: "Foo"})
	actual, err := encoder.Encode(xml)
	if err != nil {
		t.Error("Unexpected error: " + err.Error() + " Got: " + fmt.Sprintf("%x", actual))
		return
	}
	assertBinEqual(t, actual, []byte{0x45, 0x02})
}

func TestEncodePrefixDictionaryElementS(t *testing.T) {
	xml := "<s:Foo>"

	encoder := NewEncoderWithStrings(map[uint32]string{0x02: "Foo"})
	actual, err := encoder.Encode(xml)
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
	i, err := writeString(&buffer, str)
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
	i, err := writeMultiByteInt31(&buffer, num)
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

func assertBinEqual(t *testing.T, actual, expected []byte) {
	if len(actual) != len(expected) {
		t.Error("length of actual " + fmt.Sprint(len(actual)) + " not equal to length of expected " + fmt.Sprint(len(expected)))
	}
	for i, b := range actual {
		if b != expected[i] {
			fmt.Println("actual", actual, "expected", expected)
			t.Error(fmt.Sprintf("%x differs from expected %x at index %d", actual, expected, i))
		}
	}
}
