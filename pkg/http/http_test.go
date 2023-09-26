package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"ponto-menos/pkg/http/httpmock"
	"ponto-menos/pkg/text/template"
	"regexp"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

var (
	READER_ERROR = io.ErrUnexpectedEOF
)

type mockReadCloser struct{}

func (m mockReadCloser) Read(p []byte) (int, error) { return 0, READER_ERROR }
func (m mockReadCloser) Close() error               { return nil }

func createObservedLogs() (observedLogs *observer.ObservedLogs) {
	loggerCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
		os.Stderr,
		zapcore.DebugLevel,
	)

	observedCore, observedLogs := observer.New(zapcore.DebugLevel)
	logger := zap.New(zapcore.NewTee(loggerCore, observedCore))

	zap.ReplaceGlobals(logger)

	return
}

func TestHttpClientDoDecorator(t *testing.T) {
	httpmock.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		time.Sleep(time.Second)
		return nil, nil
	}

	observedLogs := createObservedLogs()

	dummyProtocolScheme := "https"
	dummyDomain := "www.dummy-domain.com.br"
	clientRequest := http.Request{}
	clientRequest.URL = &url.URL{Scheme: dummyProtocolScheme, Host: dummyDomain}
	client := &PontoMenosHTTPClient{&httpmock.MockClient{}}
	client.Do(&clientRequest)

	entry := observedLogs.All()[0]

	want := regexp.MustCompile(`^(Request took )([^:\/\s]+)( for the following url: )('[^\s]+')$`)
	got := entry.Message
	if !want.MatchString(got) {
		t.Error(template.Unexpected("message structure", got, REQUEST_TOOK_MESSAGE_TEMPLATE))
	}

	token := strings.Split(entry.Message, " ")
	time := token[2]
	url := token[7]

	wantTime := "1"
	gotTime := string(time[0])
	if gotTime != wantTime {
		t.Error(template.Unexpected("time", gotTime, wantTime))
	}

	wantUrl := fmt.Sprintf("'%s://%s'", dummyProtocolScheme, dummyDomain)
	gotUrl := url
	if gotUrl != wantUrl {
		t.Error(template.Unexpected("time", gotUrl, wantUrl))
	}
}

func TestNewHttpWrapper(t *testing.T) {
	stubResponse := &http.Response{}
	stubError := fmt.Errorf("Dummy Error")

	wrapper := NewHttpWrapper(stubResponse, stubError)

	if wrapper.Response == nil {
		template.UnexpectedValue(nil, stubResponse)
	}

	if wrapper.Error == nil {
		template.UnexpectedValue(nil, stubError)
	}
}

func TestIs2xxSuccess(t *testing.T) {
	stubResponse := &http.Response{StatusCode: 200}

	wrapper := NewHttpWrapper(stubResponse, nil)

	want := true
	got := wrapper.Is2xx()
	if got != want {
		t.Error(template.UnexpectedValue(got, want))
	}
}

func TestIs2xxFail(t *testing.T) {
	type cases struct {
		Description string
		Got         *PontoMenosHTTPResultWrapper
		Want        bool
	}

	for _, scenario := range []cases{
		{
			Description: "StatusCode 5xx",
			Got:         NewHttpWrapper(&http.Response{StatusCode: 500}, nil),
			Want:        false,
		},
		{
			Description: "StatusCode 4xx",
			Got:         NewHttpWrapper(&http.Response{StatusCode: 400}, nil),
			Want:        false,
		},
		{
			Description: "StatusCode 3xx",
			Got:         NewHttpWrapper(&http.Response{StatusCode: 300}, nil),
			Want:        false,
		},
		{
			Description: "StatusCode 1xx",
			Got:         NewHttpWrapper(&http.Response{StatusCode: 100}, nil),
			Want:        false,
		},
	} {
		t.Run(scenario.Description, func(t *testing.T) {
			if scenario.Got.Is2xx() != scenario.Want {
				t.Error(template.UnexpectedValue(scenario.Got.Is2xx(), scenario.Want))
			}
		})
	}
}

func TestResponseBody(t *testing.T) {
	dummyBody := "DummyBody"
	r := &http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte(dummyBody)))}
	wrapper := NewHttpWrapper(r, nil)

	got := wrapper.ResponseBody().Bytes()
	want := []byte(dummyBody)

	if !bytes.Equal(got, want) {
		t.Error(template.UnexpectedValue(got, want))
	}
}

func TestResponseBodyFail(t *testing.T) {
	observedLogs := createObservedLogs()

	r := &http.Response{Body: &mockReadCloser{}}
	wrapper := NewHttpWrapper(r, nil)

	if wrapper.ResponseBody() != nil {
		t.Error(template.UnexpectedValue(wrapper.ResponseBody(), nil))
	}

	entry := observedLogs.All()[0]

	want := regexp.MustCompile(`^(Could not get response body error:)(\s+)(.*)$`)
	gotMessage := entry.Message
	if !want.MatchString(gotMessage) {
		t.Error(template.Unexpected("message structure", gotMessage, RESPONSE_BODY_IO_COPY_MESSAGE_TEMPLATE))
	}

	token := strings.Split(gotMessage, ":")
	errorReason := token[1]

	if !strings.Contains(errorReason, READER_ERROR.Error()) {
		t.Error(template.Unexpected("error", errorReason, READER_ERROR.Error()))
	}
}
