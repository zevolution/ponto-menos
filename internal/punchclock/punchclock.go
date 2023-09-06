package punchclock

import (
	"encoding/json"
	"fmt"
	"ponto-menos/internal/session"
	"ponto-menos/internal/sign"
	"ponto-menos/internal/timecard"

	"go.uber.org/zap"
)

type PunchClock struct {
	Sign     sign.Signer
	Session  session.Session
	TimeCard timecard.TimeCarder
}

func (pc *PunchClock) ClockIn(credentials sign.Credentials) error {
	signedUser, err := pc.Sign.In(credentials)
	if err != nil {
		return err
	}

	b, err := json.Marshal(signedUser)
	zap.L().Debug(fmt.Sprintf("SignedUser: %v\n", string(b)))

	signedUserSession, err := pc.Session.GetUserSession(*signedUser)
	if err != nil {
		return err
	}

	b, err = json.Marshal(signedUserSession)
	zap.L().Debug(fmt.Sprintf("Session: %v\n", string(b)))

	registeredTimeCard, err := pc.TimeCard.Register(*signedUserSession)
	if err != nil {
		return err
	}

	b, err = json.Marshal(registeredTimeCard)
	zap.L().Debug(fmt.Sprintf("Registered Time Card: %v\n", string(b)))

	return nil
}
