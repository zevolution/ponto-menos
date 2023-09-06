package main

import (
	"fmt"
	"os"
	"ponto-menos/internal/punchclock"
	"ponto-menos/internal/session"
	"ponto-menos/internal/sign"
	"ponto-menos/internal/timecard"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"go.uber.org/zap"
)

func init() {
	logger := zap.Must(zap.NewDevelopment())
	if os.Getenv("APP_ENV") == "production" {
		logger = zap.Must(zap.NewProduction())
	}

	zap.ReplaceGlobals(logger)
}

func Handler(request events.CloudWatchEvent) error {
	zap.L().Debug(fmt.Sprintf("Received event of type: %v\n", request))

	user := os.Getenv("PONTO_MAIS_LOGIN")
	password := os.Getenv("PONTO_MAIS_PASSWORD")

	credentials := sign.Credentials{User: user, Password: password}

	puncher := punchclock.PunchClock{Sign: &sign.UserSign{}, Session: &session.UserSession{}, TimeCard: &timecard.TimeCard{}}

	err := puncher.ClockIn(credentials)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
