package httputil

func Is2xx(httpStatusCode int) bool {
	return httpStatusCode/100 == 2
}
