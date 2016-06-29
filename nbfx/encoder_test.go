package nbfx

import (
	"testing"
	//"io/ioutil"
	"fmt"
)

//func TestEncodeExample1(t *testing.T) {
//	encoder := NewEncoder()
//	path := "../examples/1"
//	bin, err := ioutil.ReadFile(path + ".xml")
//	if failOn(err, "unable to open " + path + ".xml", t) { return }
//	_, err = encoder.Encode(string(bin))
//	if err == nil {
//		t.Error("Expected err")
//		return
//	} else if err.Error() != "Unknown Record ID 0x44" {
//		t.Error("Expected Unknown Record ID 0x44 message but got " + err.Error())
//		return
//	}
//}

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

func assertBinEqual(t *testing.T, actual, expected []byte) {
	if len(actual) != len(expected) {
		t.Error("length of actual " + fmt.Sprint(len(actual)) + " not equal to length of expected " + fmt.Sprint(len(expected)))
	}
	for i, b := range actual {
		if b != expected[i] {
			t.Error(fmt.Sprintf("%x differs from expected %x at index %d", actual, expected, i))
		}
	}
}
