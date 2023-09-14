package timecard

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	gohttp "net/http"
	"ponto-menos/internal/session"
	"ponto-menos/internal/sign"
	"ponto-menos/pkg/http"
	"ponto-menos/pkg/http/httperror"
	"ponto-menos/pkg/http/httpmock"
	"ponto-menos/pkg/text/template"
	"reflect"
	"strconv"
	"testing"
)

var (
	dummyData = struct {
		ApiVersion   int
		Success      string
		Token        string
		ClientId     string
		Login        string
		SignInCount  int
		LastSignInIp string
		LastSignInAt int
		SessionToken string
		SessionUuid  string
	}{
		ApiVersion:   session.API_VERSION,
		Success:      "Login efetuado com sucesso!",
		Token:        "token",
		ClientId:     "clientId",
		Login:        "user@user.com.br",
		SignInCount:  1285,
		LastSignInIp: "172.59.28.255",
		LastSignInAt: 1692815760,
		SessionToken: "cgS1fidp0dsq1puS_BTxS3pM4spAvG-YcPongfaD4HddabiaIdskdbb",
		SessionUuid:  "bbfce251-a30b-4b13-88e8-2b331e2a37d9",
	}
)

func init() {
	http.Client = &httpmock.MockClient{}
}

func TestError(t *testing.T) {
	want := "Simulate http client error"

	httpmock.GetDoFunc = func(req *gohttp.Request) (*gohttp.Response, error) {
		return nil, fmt.Errorf(want)
	}

	stubSignedUserDetail := &sign.SignedUserDetail{Login: dummyData.Login, SignInCount: dummyData.SignInCount, LastSignInIp: dummyData.LastSignInIp, LastSignInAt: dummyData.LastSignInAt}
	stubSignedUser := &sign.SignedUser{Success: dummyData.Success, Token: dummyData.Token, ClientId: dummyData.ClientId, Detail: stubSignedUserDetail}
	signedUserSession := &sign.SignedUserSession{ApiVersion: dummyData.ApiVersion, Uid: dummyData.Login, Uuid: dummyData.SessionUuid, AuthorizationToken: dummyData.SessionToken, SignedUser: stubSignedUser}

	sut := &TimeCard{}
	_, got := sut.Register(*signedUserSession)

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

	stubSignedUserDetail := &sign.SignedUserDetail{Login: dummyData.Login, SignInCount: dummyData.SignInCount, LastSignInIp: dummyData.LastSignInIp, LastSignInAt: dummyData.LastSignInAt}
	stubSignedUser := &sign.SignedUser{Success: dummyData.Success, Token: dummyData.Token, ClientId: dummyData.ClientId, Detail: stubSignedUserDetail}
	signedUserSession := &sign.SignedUserSession{ApiVersion: dummyData.ApiVersion, Uid: dummyData.Login, Uuid: dummyData.SessionUuid, AuthorizationToken: dummyData.SessionToken, SignedUser: stubSignedUser}

	sut := &TimeCard{}
	_, got := sut.Register(*signedUserSession)

	var want httperror.ErrDiffFrom2xx
	if !errors.As(got, &want) {
		t.Errorf(template.Unexpected("header error instance", reflect.TypeOf(got), reflect.TypeOf(want)))
	}
}

func TestCreateTimeCard(t *testing.T) {
	responseFile, _ := ioutil.ReadFile("testdata/timecard_response.json")
	response := ioutil.NopCloser(bytes.NewReader(responseFile))

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
		SessionUuid  string
	}{
		ApiVersion:   session.API_VERSION,
		Success:      "Login efetuado com sucesso!",
		Token:        "token",
		ClientId:     "clientId",
		Login:        "user@user.com.br",
		SignInCount:  1285,
		LastSignInIp: "172.59.28.255",
		LastSignInAt: 1692815760,
		SessionToken: "cgS1fidp0dsq1puS_BTxS3pM4spAvG-YcPongfaD4HddabiaIdskdbb",
		SessionUuid:  "bbfce251-a30b-4b13-88e8-2b331e2a37d9",
	}

	httpmock.GetDoFunc = func(r *gohttp.Request) (*gohttp.Response, error) {
		b, _ := io.ReadAll(r.Body)
		got := make(map[string]*interface{})
		json.Unmarshal(b, &got)
		fmt.Println(got)

		requestFile, _ := ioutil.ReadFile("testdata/timecard_request.json")
		want := make(map[string]*interface{})
		json.Unmarshal(requestFile, &want)
		fmt.Println(want)

		if !reflect.DeepEqual(got, want) {
			t.Errorf(template.Unexpected("request body", got, want))
		}

		headers := map[string]string{
			"api-version":  strconv.Itoa(session.API_VERSION),
			"access-token": dummyData.Token,
			"client":       dummyData.ClientId,
			"origin":       "https://app2.pontomais.com.br",
			"referer":      "https://app2.pontomais.com.br",
			"token":        dummyData.Token,
			"uid":          dummyData.Login,
			"uuid":         dummyData.SessionUuid,
			"user-agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.115 Safari/537.36",
		}
		for k, v := range headers {
			got := r.Header.Get(k)
			if r.Header.Get(k) != v {
				t.Errorf(template.Unexpected("header value", got, v))
			}
		}

		httpMethod := "POST"
		if r.Method != httpMethod {
			t.Errorf(template.Unexpected("http method", httpMethod, r.Method))
		}

		sessionUrlPath := "/api/time_cards/register"
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
	signedUserSession := &sign.SignedUserSession{ApiVersion: dummyData.ApiVersion, Uid: dummyData.Login, Uuid: dummyData.SessionUuid, AuthorizationToken: dummyData.SessionToken, SignedUser: stubSignedUser}

	sut := &TimeCard{}
	got, _ := sut.Register(*signedUserSession)

	var want *RegisteredTimeCard
	json.Unmarshal(responseFile, &want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf(template.Unexpected("response", got, want))
	}
}
