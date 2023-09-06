package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"ponto-menos/pkg/http/httputil"
	"time"

	"go.uber.org/zap"
)

type (
	HTTPClient interface {
		Do(req *http.Request) (*http.Response, error)
	}

	PontoMenosHTTPClient struct {
		HttpClient http.Client
	}

	PontoMenosHTTPResultWrapper struct {
		Response *http.Response
		Error    error
	}
)

var (
	Client HTTPClient
)

func init() {
	Client = &PontoMenosHTTPClient{http.Client{}}
}

func (pmhc *PontoMenosHTTPClient) Do(req *http.Request) (res *http.Response, err error) {
	start := time.Now()
	res, err = pmhc.HttpClient.Do(req)
	elapsed := time.Since(start)

	zap.L().Debug(fmt.Sprintf("Request took %s for the following url: '%v'", elapsed, req.URL))

	return
}

func NewHttpWrapper(res *http.Response, err error) *PontoMenosHTTPResultWrapper {
	return &PontoMenosHTTPResultWrapper{
		Response: res,
		Error:    err,
	}
}

func (rw *PontoMenosHTTPResultWrapper) Is2xx() bool {
	return httputil.Is2xx(rw.Response.StatusCode)
}

func (rw *PontoMenosHTTPResultWrapper) ResponseBody() *bytes.Buffer {
	defer rw.Response.Body.Close()

	var b bytes.Buffer
	if _, err := io.Copy(&b, rw.Response.Body); err != nil {
		zap.L().Error(fmt.Sprintf("Could not get response body error: %v", err.Error()))
		return nil
	}

	return &b
}
