package httperror

import (
	"bytes"
	"errors"
	"io/ioutil"
	gohttp "net/http"
	"ponto-menos/pkg/http"
	"ponto-menos/pkg/text/template"
	"regexp"
	"testing"
)

func TestNew(t *testing.T) {
	rw := http.PontoMenosHTTPResultWrapper{Response: nil, Error: errors.New("teste")}
	sut := NewErrDiffFrom2xx(rw)

	var edf2 ErrDiffFrom2xx
	if !errors.As(sut, &edf2) {
		t.Error("Diff instance")
	}
}

func TestErrorMessage(t *testing.T) {
	dummyResponseBodyAsString := `{"keyTest": "valueTest"}`
	responseBody := ioutil.NopCloser(bytes.NewReader([]byte(dummyResponseBodyAsString)))
	stubResponse := gohttp.Response{StatusCode: 500, Body: responseBody}
	rw := http.PontoMenosHTTPResultWrapper{Response: &stubResponse, Error: errors.New("teste")}

	sut := NewErrDiffFrom2xx(rw)

	got := sut.Error()
	want := regexp.MustCompile(`^(Unsuccessful request for the following reason -> HTTP Code: )([^\s]+)( \| )(ResponseBody: )(.*)$`)
	if !want.MatchString(got) {
		t.Error(template.Unexpected("message structure", got, UNSUCCESSFUL_REQUEST_TEMPLATE_ERROR))
	}
}
