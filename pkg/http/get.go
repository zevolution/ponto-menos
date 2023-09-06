package http

import (
	"log"
	"net/http"
)

func DoHTTPGetWithHeaders(requestURL string, headers map[string]string) *PontoMenosHTTPResultWrapper {
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
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
