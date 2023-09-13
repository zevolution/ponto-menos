package session

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	gohttp "net/http"
	"ponto-menos/internal/sign"
	"ponto-menos/pkg/http"
	"ponto-menos/pkg/http/httperror"
	"ponto-menos/pkg/http/httpmock"
	"ponto-menos/pkg/text/template"
	"reflect"
	"strconv"
	"testing"
)

func init() {
	http.Client = &httpmock.MockClient{}
}

func TestError(t *testing.T) {
	want := "Simulate http client error"

	httpmock.GetDoFunc = func(req *gohttp.Request) (*gohttp.Response, error) {
		return nil, fmt.Errorf(want)
	}

	stubSignedUserDetail := &sign.SignedUserDetail{Login: "dummyLogin", SignInCount: 0, LastSignInIp: "dummyIP", LastSignInAt: 0}
	stubSignedUser := &sign.SignedUser{Success: "dummyOk", Token: "dummyToken", ClientId: "dummyClientId", Detail: stubSignedUserDetail}

	sut := &UserSession{}
	_, got := sut.GetUserSession(*stubSignedUser)

	if got.Error() != want {
		t.Errorf(template.Unexpected("header error message", got.Error(), want))
	}
}

func TestErrorDifferentFrom2xx(t *testing.T) {
	httpmock.GetDoFunc = func(req *gohttp.Request) (*gohttp.Response, error) {
		return &gohttp.Response{
			StatusCode: 401,
			Body:       nil,
		}, nil
	}

	stubSignedUserDetail := &sign.SignedUserDetail{Login: "dummyLogin", SignInCount: 0, LastSignInIp: "dummyIP", LastSignInAt: 0}
	stubSignedUser := &sign.SignedUser{Success: "dummyOk", Token: "dummyToken", ClientId: "dummyClientId", Detail: stubSignedUserDetail}

	sut := &UserSession{}
	_, got := sut.GetUserSession(*stubSignedUser)

	var want httperror.ErrDiffFrom2xx
	if !errors.As(got, &want) {
		t.Errorf(template.Unexpected("header error instance", reflect.TypeOf(got), reflect.TypeOf(want)))
	}
}

func TestGetUserSession(t *testing.T) {
	file, _ := ioutil.ReadFile("testdata/session_response.json")
	response := ioutil.NopCloser(bytes.NewReader(file))

	var sessionResponse SessionResponseWrapper
	json.Unmarshal(file, &sessionResponse)

	dummyData := struct {
		ApiVersion   int
		Success      string
		Token        string
		ClientId     string
		Login        string
		SignInCount  int
		LastSignInIp string
		LastSignInAt int
		SessionToken string
	}{
		ApiVersion:   API_VERSION,
		Success:      "dummyOk",
		Token:        "dummyToken",
		ClientId:     "dummyClientId",
		Login:        "dummyLogin",
		SignInCount:  0,
		LastSignInIp: "dummyIp",
		LastSignInAt: 0,
		SessionToken: sessionResponse.Session.Client.Token,
	}

	httpmock.GetDoFunc = func(r *gohttp.Request) (*gohttp.Response, error) {
		headers := map[string]string{
			"api-version":  strconv.Itoa(dummyData.ApiVersion),
			"uid":          dummyData.Login,
			"uuid":         SessionUuid,
			"client":       dummyData.ClientId,
			"access-token": dummyData.Token,
			"token":        dummyData.Token,
		}
		for k, v := range headers {
			got := r.Header.Get(k)
			if r.Header.Get(k) != v {
				t.Errorf(template.Unexpected("header value", got, v))
			}
		}

		httpMethod := "GET"
		if r.Method != httpMethod {
			t.Errorf(template.Unexpected("http method", httpMethod, r.Method))
		}

		sessionUrlPath := "/api/session"
		if r.URL.Path != sessionUrlPath {
			t.Errorf(template.Unexpected("path", sessionUrlPath, r.URL.Path))
		}

		return &gohttp.Response{
			StatusCode: 200,
			Body:       response,
		}, nil
	}

	stubSignedUserDetail := &sign.SignedUserDetail{Login: dummyData.Login, SignInCount: dummyData.SignInCount, LastSignInIp: dummyData.LastSignInIp, LastSignInAt: dummyData.LastSignInAt}
	stubSignedUser := &sign.SignedUser{Success: dummyData.Success, Token: dummyData.Token, ClientId: dummyData.ClientId, Detail: stubSignedUserDetail}

	sut := &UserSession{}
	got, _ := sut.GetUserSession(*stubSignedUser)

	want := &sign.SignedUserSession{ApiVersion: dummyData.ApiVersion, Uid: dummyData.Login, Uuid: SessionUuid, AuthorizationToken: dummyData.SessionToken, SignedUser: stubSignedUser}
	if !reflect.DeepEqual(got, want) {
		t.Errorf(template.Unexpected("object", got, want))
	}
}
