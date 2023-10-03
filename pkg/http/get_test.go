package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"ponto-menos/pkg/http/httpmock"
	"ponto-menos/pkg/text/template"
	"testing"
)

func init() {
	Client = &httpmock.MockClient{}
}

func TestDoHTTPGetWithHeaders(t *testing.T) {
	dummyResponseBodyAsString := `{"keyTest": "valueTest"}`
	mockResponse := ioutil.NopCloser(bytes.NewReader([]byte(dummyResponseBodyAsString)))

	stubHttpStatusCode := 200
	stubRequestUrl := "http://localhost:3000/test"
	stubHeaders := map[string]string{
		"dummyKeyOne": "dummyValueOne",
		"dummyKeyTwo": "dummyValueTwo",
	}

	httpmock.GetDoFunc = func(r *http.Request) (*http.Response, error) {
		for k, v := range stubHeaders {
			got := r.Header.Get(k)
			if r.Header.Get(k) != v {
				t.Errorf(template.Unexpected("header value", got, v))
			}
		}

		httpMethod := "GET"
		if r.Method != httpMethod {
			t.Errorf(template.Unexpected("http method", r.Method, httpMethod))
		}

		sessionUrlPath := "/test"
		if r.URL.Path != sessionUrlPath {
			t.Errorf(template.Unexpected("path", r.URL.Path, sessionUrlPath))
		}

		return &http.Response{
			StatusCode: stubHttpStatusCode,
			Body:       mockResponse,
		}, nil
	}

	rw := DoHTTPGetWithHeaders(stubRequestUrl, stubHeaders)

	if !rw.Is2xx() {
		t.Errorf(template.Unexpected("http status code", rw.Response.StatusCode, stubHttpStatusCode))
	}

	responseBody := rw.ResponseBody().Bytes()
	got := string(responseBody)
	if got != dummyResponseBodyAsString {
		t.Errorf(template.Unexpected("response body", got, dummyResponseBodyAsString))
	}

}

func TestDoHTTPGetWithHeadersErrors(t *testing.T) {
	mockError := fmt.Errorf("Some error")

	httpmock.GetDoFunc = func(r *http.Request) (*http.Response, error) {
		return nil, mockError
	}

	rw := DoHTTPGetWithHeaders("", nil)

	if rw.Error == nil {
		t.Errorf(template.Unexpected("error", rw.Error, mockError))
	}

}

func TestDoHTTPGetWithDefaultHeaders(t *testing.T) {
	httpmock.GetDoFunc = func(r *http.Request) (*http.Response, error) {
		want := "application/json"
		got := r.Header.Get("Content-Type")

		if got != want {
			t.Errorf(template.Unexpected("default header", got, want))
		}

		return &http.Response{
			StatusCode: 500,
			Body:       nil,
		}, nil
	}

	_ = DoHTTPGetWithHeaders("", nil)
}
