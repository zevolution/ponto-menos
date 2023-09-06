package session

import (
	"encoding/json"
	"ponto-menos/internal/sign"
	"ponto-menos/pkg/http"
	"ponto-menos/pkg/http/httperror"
	"strconv"

	"github.com/google/uuid"
)

type (
	Session interface {
		GetUserSession(signedUser sign.SignedUser) (*sign.SignedUserSession, error)
	}

	UserSession struct {
	}

	SessionResponse struct {
		Client *SessionClientResponse `json:"client"`
	}

	SessionClientResponse struct {
		Token string `json:"time_clocks_token"`
	}

	SessionResponseWrapper struct {
		Session SessionResponse `json:"session"`
	}
)

const (
	API_VERSION = 2
)

var (
	SessionUuid string
)

func init() {
	SessionUuid = uuid.New().String()
}

func (us *UserSession) GetUserSession(signedUser sign.SignedUser) (*sign.SignedUserSession, error) {
	headers := headersFrom(signedUser)

	rw := http.DoHTTPGetWithHeaders("https://api.pontomais.com.br/api/session", headers)

	if rw.Error != nil {
		return nil, rw.Error
	} else if !rw.Is2xx() {
		return nil, httperror.NewErrDiffFrom2xx(*rw)
	}

	var response SessionResponseWrapper
	json.Unmarshal(rw.ResponseBody().Bytes(), &response)

	userSession := sign.SignedUserSession{
		ApiVersion:         API_VERSION,
		Uid:                signedUser.Detail.Login,
		Uuid:               SessionUuid,
		AuthorizationToken: response.Session.Client.Token,
		SignedUser:         &signedUser,
	}

	return &userSession, nil
}

func headersFrom(signedUser sign.SignedUser) map[string]string {
	return map[string]string{
		"api-version":  strconv.Itoa(API_VERSION),
		"uid":          signedUser.Detail.Login,
		"uuid":         SessionUuid,
		"client":       signedUser.ClientId,
		"access-token": signedUser.Token,
		"token":        signedUser.Token,
	}
}
