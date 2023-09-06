package http

import (
	"bytes"
	"log"
	"net/http"
)

func DoHTTPPostWithBody(requestURL, requestBody string) *PontoMenosHTTPResultWrapper {
	return DoHTTPPostWithBodyAndHeaders(requestURL, requestBody, nil)
}

func DoHTTPPostWithBodyAndHeaders(requestURL, requestBody string, headers map[string]string) *PontoMenosHTTPResultWrapper {
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		log.Panicf("client: could not create request: %s\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	for key := range headers {
		req.Header.Set(key, headers[key])
	}

	response, err := Client.Do(req)

	if err != nil {
		return NewHttpWrapper(nil, err)
	}

	return NewHttpWrapper(response, nil)
}
