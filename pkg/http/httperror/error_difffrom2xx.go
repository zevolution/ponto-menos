package httperror

import (
	"fmt"
	"ponto-menos/pkg/http"
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
	return fmt.Sprintf("Unsuccessful request for the following reason -> HTTP Code: %d | ResponseBody: %v", err.pmh.Response.StatusCode, string(err.pmh.ResponseBody().Bytes()))
}
