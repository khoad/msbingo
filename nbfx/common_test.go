package nbfx

import (
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, actual, expected interface{}) {
	if expected != actual {
		t.Errorf("%v not equal to expected %v", actual, expected)
	}
}

func assertStringEqual(t *testing.T, actual, expected string) {
	if expected != actual {
		t.Error(actual + " not equal to expected " + expected)
	}
}

func assertBinEqual(t *testing.T, actual, expected []byte) {
	if len(actual) != len(expected) {
		t.Error("length of actual " + fmt.Sprint(len(actual)) + " not equal to length of expected " + fmt.Sprint(len(expected)))
	}
	for i, b := range actual {
		if i == len(expected) || b != expected[i] {
			pointerLine(t, actual, expected, i)
			return
		}
	}
	if len(actual) != len(expected) {
		pointerLine(t, actual, expected, len(actual))
	}
}

func pointerLine(t *testing.T, actual, expected []byte, i int) {
	pointerLine := "^^"
	for j := 1; j <= i; j++ {
		pointerLine = "--" + pointerLine
	}
	t.Error(fmt.Sprintf("actual\n%x\ndiffers from expected at index %d\n%x\n%s\n", actual, i, expected, pointerLine))
}
