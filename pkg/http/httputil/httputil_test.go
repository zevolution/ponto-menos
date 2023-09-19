package httputil

import (
	"ponto-menos/pkg/text/template"
	"testing"
)

func TestSuccessIs2xx(t *testing.T) {
	want := true
	got := Is2xx(200)

	if got != want {
		t.Error(template.UnexpectedValue(got, want))
	}
}

func TestFailIs2xx(t *testing.T) {
	type cases struct {
		Description    string
		HttpStatusCode int
		Want           bool
	}

	for _, scenario := range []cases{
		{
			Description:    "StatusCode 5xx",
			HttpStatusCode: 500,
			Want:           false,
		},
		{
			Description:    "StatusCode 4xx",
			HttpStatusCode: 400,
			Want:           false,
		},
		{
			Description:    "StatusCode 3xx",
			HttpStatusCode: 300,
			Want:           false,
		},
		{
			Description:    "StatusCode 1xx",
			HttpStatusCode: 100,
			Want:           false,
		},
	} {
		t.Run(scenario.Description, func(t *testing.T) {
			got := Is2xx(scenario.HttpStatusCode)
			if got != scenario.Want {
				t.Error(template.UnexpectedValue(got, scenario.Want))
			}
		})
	}
}
