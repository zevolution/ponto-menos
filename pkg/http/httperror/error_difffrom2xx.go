package httperror

import (
	"fmt"
	"ponto-menos/pkg/http"
)

const (
	UNSUCCESSFUL_REQUEST_TEMPLATE_ERROR = "Unsuccessful request for the following reason -> HTTP Code: %d | ResponseBody: %v"
)

type ErrDiffFrom2xx struct {
	pmh http.PontoMenosHTTPResultWrapper
}

func NewErrDiffFrom2xx(pmh http.PontoMenosHTTPResultWrapper) error {
	return ErrDiffFrom2xx{
		pmh: pmh,
	}
}

func (err ErrDiffFrom2xx) Error() string {
	return fmt.Sprintf(UNSUCCESSFUL_REQUEST_TEMPLATE_ERROR, err.pmh.Response.StatusCode, string(err.pmh.ResponseBody().Bytes()))
}
