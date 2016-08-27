package nbfs

import (
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, actual, expected string) {
	if expected != actual {
		t.Error(actual + "\nnot equal to expected\n" + expected)
	}
}

func assertBinEqual(t *testing.T, actual, expected []byte) {
	if len(actual) != len(expected) {
		t.Error("length of actual " + fmt.Sprint(len(actual)) + " not equal to length of expected " + fmt.Sprint(len(expected)))
	}
	for i, b := range actual {
		if b != expected[i] {
			pointerLine := "^"
			for j := 1; j <= i; j++ {
				pointerLine = "--" + pointerLine
			}
			t.Error(fmt.Sprintf("actual\n%x\ndiffers from expected at index %d\n%x\n%s\n", actual, i, expected, pointerLine))
			return
		}
	}
}

func failOn(err error, message string, t *testing.T) bool {
	if err != nil {
		t.Error(message)
		return true
	}
	return false
}
