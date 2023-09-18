package template

import (
	"fmt"
	"strings"
	"testing"
)

func TestUnexpectedTemplate(t *testing.T) {
	dummyReason := "input"
	dummyWant := 1
	dummyGot := 2

	got := Unexpected(dummyReason, dummyWant, dummyGot)

	want := fmt.Sprintf("Unexpected %v.\nwant: %v\ngot: %v", dummyReason, dummyGot, dummyWant)
	if got != want {
		t.Errorf("'Unexpected' template is not equals.\nwant: %v\ngot: %v", strings.Replace(want, "\n", " ", -1), strings.Replace(got, "\n", " ", -1))
	}
}

func TestUnexpectedValueTemplate(t *testing.T) {
	dummyWant := 1
	dummyGot := 2

	got := UnexpectedValue(dummyWant, dummyGot)

	want := fmt.Sprintf("Unexpected value.\nwant: %v\ngot: %v", dummyGot, dummyWant)
	if got != want {
		t.Errorf("'UnexpectedValue' template is not equals.\nwant: %v\ngot: %v", strings.Replace(want, "\n", " ", -1), strings.Replace(got, "\n", " ", -1))
	}
}
