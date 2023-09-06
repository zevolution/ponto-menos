package template

import "fmt"

func Unexpected(reason string, want any, got any) string {
	return fmt.Sprintf("Unexpected %v.\nwant: %v\ngot: %v", reason, want, got)
}

func UnexpectedValue(want any, got any) string {
	return Unexpected("value", want, got)
}
