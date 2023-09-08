package template

import "fmt"

func Unexpected(reason string, got any, want any) string {
	return fmt.Sprintf("Unexpected %v.\nwant: %v\ngot: %v", reason, want, got)
}

func UnexpectedValue(got any, want any) string {
	return Unexpected("value", got, want)
}
