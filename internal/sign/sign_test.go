package sign

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	gohttp "net/http"
	"ponto-menos/pkg/http"
	"ponto-menos/pkg/http/httperror"
	"ponto-menos/pkg/http/httpmock"
	"ponto-menos/pkg/text/template"
	"reflect"
	"strings"
	"testing"
)

func init() {
	http.Client = &httpmock.MockClient{}
}

func TestError(t *testing.T) {
	dummyUser := "user"
	dummyPassword := "password"
	want := "Simulate http client error"

	httpmock.GetDoFunc = func(r *gohttp.Request) (*gohttp.Response, error) {
		return nil, fmt.Errorf(want)
	}

	stubCredentials := &Credentials{User: dummyUser, Password: dummyPassword}
	sut := &UserSign{}
	_, got := sut.In(*stubCredentials)

	if got.Error() != want {
		t.Errorf(template.Unexpected("error message", got.Error(), want))
	}
}

func TestErrorDifferentFrom2xx(t *testing.T) {
	dummyUser := "user"
	dummyPassword := "password"
	httpmock.GetDoFunc = func(r *gohttp.Request) (*gohttp.Response, error) {
		return &gohttp.Response{
			StatusCode: 401,
			Body:       nil,
		}, nil
	}

	stubCredentials := &Credentials{User: dummyUser, Password: dummyPassword}
	sut := &UserSign{}
	_, got := sut.In(*stubCredentials)

	var want httperror.ErrDiffFrom2xx
	if !errors.As(got, &want) {
		t.Errorf(template.Unexpected("error instance", reflect.TypeOf(got), reflect.TypeOf(want)))
	}
}

func TestUserSigninWithCredentials(t *testing.T) {
	file, _ := ioutil.ReadFile("testdata/signin_response.json")
	response := ioutil.NopCloser(bytes.NewReader(file))

	var want SignedUser
	json.Unmarshal(file, &want)

	dummyUser := "user"
	dummyPassword := "password"
	httpmock.GetDoFunc = func(r *gohttp.Request) (*gohttp.Response, error) {
		var c Credentials
		json.NewDecoder(r.Body).Decode(&c)

		if dummyUser != c.User {
			t.Errorf(template.UnexpectedValue(dummyUser, c.User))
		}

		if dummyPassword != c.Password {
			t.Errorf(template.UnexpectedValue(dummyPassword, c.Password))
		}

		httpMethod := "POST"
		if r.Method != httpMethod {
			t.Errorf(template.UnexpectedValue(httpMethod, r.Method))
		}

		signinUrl := "/api/auth/sign_in"
		if !strings.Contains(r.URL.Path, signinUrl) {
			t.Errorf(template.UnexpectedValue(signinUrl, r.URL.Path))
		}

		return &gohttp.Response{
			StatusCode: 201,
			Body:       response,
		}, nil
	}

	stubCredentials := &Credentials{User: dummyUser, Password: dummyPassword}
	sut := &UserSign{}
	got, _ := sut.In(*stubCredentials)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf(template.Unexpected("object", got, want))
	}
}
