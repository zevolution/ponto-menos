package punchclock

import (
	"fmt"
	"ponto-menos/internal/session"
	"ponto-menos/internal/sign"
	"ponto-menos/internal/timecard"
	"ponto-menos/pkg/text/template"
	"testing"
)

const (
	TIMES_TO_CALL int = 1
)

var (
	dummyData = struct {
		User         string
		Password     string
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
		User:         "dummy@user.com.br",
		Password:     "DuMMyP455W0RD",
		ApiVersion:   session.API_VERSION,
		Success:      "Login efetuado com sucesso!",
		Token:        "token",
		ClientId:     "clientId",
		Login:        "dummy@user.com.br",
		SignInCount:  1285,
		LastSignInIp: "172.59.28.255",
		LastSignInAt: 1692815760,
		SessionToken: "cgS1fidp0dsq1puS_BTxS3pM4spAvG-YcPongfaD4HddabiaIdskdbb",
		SessionUuid:  "bbfce251-a30b-4b13-88e8-2b331e2a37d9",
	}
	dummyCredentials = &sign.Credentials{User: dummyData.User, Password: dummyData.Password}
	stubSignedUser   = &sign.SignedUser{
		Success:  dummyData.Success,
		Token:    dummyData.Token,
		ClientId: dummyData.ClientId,
		Detail: &sign.SignedUserDetail{
			Login:        dummyData.User,
			SignInCount:  dummyData.SignInCount,
			LastSignInIp: dummyData.LastSignInIp,
			LastSignInAt: dummyData.LastSignInAt},
	}
	stubSignedUserSession = &sign.SignedUserSession{
		ApiVersion:         dummyData.ApiVersion,
		Uid:                dummyData.Login,
		Uuid:               dummyData.SessionUuid,
		AuthorizationToken: dummyData.SessionToken,
		SignedUser:         stubSignedUser,
	}
)

type (
	UserSignSpy struct {
		times int
	}

	UserSessionSpy struct {
		times int
	}

	TimeCardSpy struct {
		times int
	}
)

func (uss *UserSignSpy) In(credentials sign.Credentials) (*sign.SignedUser, error) {
	uss.times++

	if credentials != *dummyCredentials {
		return nil, fmt.Errorf("Credentials input is different from stubCredentials")
	}

	return stubSignedUser, nil
}

func (uss *UserSessionSpy) GetUserSession(signedUser sign.SignedUser) (*sign.SignedUserSession, error) {
	uss.times++

	if signedUser != *stubSignedUser {
		return nil, fmt.Errorf("SignedUser input is different from stubSignedUser")
	}

	return stubSignedUserSession, nil
}

func (tcs *TimeCardSpy) Register(signedUserSession sign.SignedUserSession) (*timecard.RegisteredTimeCard, error) {
	tcs.times++

	if signedUserSession != *stubSignedUserSession {
		return nil, fmt.Errorf("SignedUserSession input is different from stubSignedUserSession")
	}

	return &timecard.RegisteredTimeCard{}, nil
}

func TestPunchClock(t *testing.T) {
	userSignSpy := &UserSignSpy{}
	userSessionSpy := &UserSessionSpy{}
	timeCardSpy := &TimeCardSpy{}
	stubCredentials := sign.Credentials{User: dummyData.User, Password: dummyData.Password}

	sut := PunchClock{Sign: userSignSpy, Session: userSessionSpy, TimeCard: timeCardSpy}
	err := sut.ClockIn(stubCredentials)

	if err != nil {
		t.Errorf(err.Error())
	}

	var want = TIMES_TO_CALL
	var got = -1

	got = userSignSpy.times
	if got != want {
		t.Errorf(template.Unexpected("'userSign' called times", got, want))
	}

	got = userSessionSpy.times
	if got != want {
		t.Errorf(template.Unexpected("'userSession' called times", got, want))
	}

	got = timeCardSpy.times
	if got != want {
		t.Errorf(template.Unexpected("'timeCard' called times", got, want))
	}

}
