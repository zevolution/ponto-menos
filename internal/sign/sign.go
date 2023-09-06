package sign

import (
	"encoding/json"
	"ponto-menos/pkg/http"
	"ponto-menos/pkg/http/httperror"
)

type (
	Signer interface {
		In(credentials Credentials) (*SignedUser, error)
	}

	UserSign struct {
	}

	Credentials struct {
		User     string `json:"login"`
		Password string `json:"password"`
	}

	SignedUser struct {
		Success  string            `json:"success"`
		Token    string            `json:"token"`
		ClientId string            `json:"client_id"`
		Detail   *SignedUserDetail `json:"data"`
	}

	SignedUserDetail struct {
		Login        string `json:"login"`
		SignInCount  int    `json:"sign_in_count"`
		LastSignInIp string `json:"last_sign_in_ip"`
		LastSignInAt int    `json:"last_sign_in_at"`
	}

	SignedUserSession struct {
		ApiVersion         int         `json:"apiVersion"`
		Uid                string      `json:"uid"`
		Uuid               string      `json:"uuid"`
		AuthorizationToken string      `json:"authorizationToken"`
		SignedUser         *SignedUser `json:"signedUser"`
	}
)

func (us *UserSign) In(credentials Credentials) (*SignedUser, error) {
	requestBody, _ := json.Marshal(credentials)

	rw := http.DoHTTPPostWithBody("https://api.pontomais.com.br/api/auth/sign_in", string(requestBody))

	if rw.Error != nil {
		return nil, rw.Error
	} else if !rw.Is2xx() {
		return nil, httperror.NewErrDiffFrom2xx(*rw)
	}

	var signedUser SignedUser
	json.Unmarshal(rw.ResponseBody().Bytes(), &signedUser)

	return &signedUser, nil
}
