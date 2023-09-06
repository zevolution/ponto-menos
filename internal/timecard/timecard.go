package timecard

import (
	"encoding/json"
	"fmt"
	"ponto-menos/internal/session"
	"ponto-menos/internal/sign"
	"ponto-menos/pkg/http"
	"ponto-menos/pkg/http/httperror"
	"strconv"
)

type (
	TimeCarder interface {
		Register(signedUserSession sign.SignedUserSession) (*RegisteredTimeCard, error)
	}

	TimeCard struct{}

	RegisteredTimeCard struct {
		Controller    string `json:"controller"`
		Action        string `json:"action"`
		Timestamp     string `json:"timestamp"`
		Uuid          string `json:"uuid"`
		CurrentUserId string `json:"current_user_id"`
		RemoteIp      string `json:"remote_ip"`
		MessageId     string `json:"message_id"`
	}

	timeRegisterRequest struct {
		Image      *interface{}                `json:"image"`
		Employee   timeRegisterEmployeeRequest `json:"employee"`
		TimeCard   timeCardRegisterRequest     `json:"time_card"`
		Path       string                      `json:"_path"`
		AppVersion string                      `json:"_appVersion"`
		Device     timeRegisterDeviceRequest   `json:"_device"`
	}

	timeRegisterEmployeeRequest struct {
		Id  *interface{} `json:"id"`
		Pin *interface{} `json:"pin"`
	}

	timeCardRegisterRequest struct {
		Latitude          int          `json:"latitude"`
		Longitude         int          `json:"longitude"`
		Address           string       `json:"address"`
		ReferenceId       *interface{} `json:"reference_id"`
		OriginalLatitude  string       `json:"original_latitude"`
		OriginalLongitude string       `json:"original_longitude"`
		OriginalAddress   string       `json:"original_address"`
		LocationEdited    bool         `json:"location_edited"`
		Accuracy          int          `json:"accuracy"`
		AccuracyMethod    *interface{} `json:"accuracy_method"`
		Image             *interface{} `json:"image"`
	}

	timeRegisterDeviceRequest struct {
		Manufacturer string                        `json:"manufacturer"`
		Model        string                        `json:"model"`
		Uuid         timeRegisterDeviceUuidRequest `json:"uuid"`
		Version      string                        `json:"version"`
	}

	timeRegisterDeviceUuidRequest struct {
		Success       string                              `json:"success"`
		Token         string                              `json:"token"`
		ClientId      string                              `json:"client_id"`
		Detail        timeRegisterDeviceUuidDetailRequest `json:"data"`
		Uuid          string                              `json:"uuid"`
		Authorization string                              `json:"authorization"`
	}

	timeRegisterDeviceUuidDetailRequest struct {
		Email        string `json:"email"`
		SignInCount  int    `json:"sign_in_count"`
		LastSignInIp string `json:"last_sign_in_ip"`
		LastSignInAt int    `json:"last_sign_in_at"`
	}

	timeCardRegisterResponse struct {
		Controller    string `json:"controller"`
		Action        string `json:"action"`
		Timestamp     string `json:"timestamp"`
		Uuid          string `json:"uuid"`
		CurrentUserId string `json:"current_user_id"`
		RemoteIp      string `json:"remote_ip"`
		MessageId     string `json:"message_id"`
	}
)

func (tc *TimeCard) Register(signedUserSession sign.SignedUserSession) (*RegisteredTimeCard, error) {
	requestBody, _ := json.Marshal(buildRequestBody(signedUserSession))

	rw := http.DoHTTPPostWithBodyAndHeaders("https://api.pontomais.com.br/api/time_cards/register", string(requestBody), headersFrom(signedUserSession))

	if rw.Error != nil {
		return nil, rw.Error
	} else if !rw.Is2xx() {
		return nil, httperror.NewErrDiffFrom2xx(*rw)
	}

	var timeCardRegisterResponse timeCardRegisterResponse
	json.Unmarshal(rw.ResponseBody().Bytes(), &timeCardRegisterResponse)

	return &RegisteredTimeCard{
		Controller:    timeCardRegisterResponse.Controller,
		Action:        timeCardRegisterResponse.Action,
		Timestamp:     timeCardRegisterResponse.Timestamp,
		Uuid:          timeCardRegisterResponse.Uuid,
		CurrentUserId: timeCardRegisterResponse.CurrentUserId,
		RemoteIp:      timeCardRegisterResponse.RemoteIp,
		MessageId:     timeCardRegisterResponse.MessageId,
	}, nil
}

func buildRequestBody(signedUserSession sign.SignedUserSession) timeRegisterRequest {
	return timeRegisterRequest{
		Path: "/registrar-ponto",
		TimeCard: timeCardRegisterRequest{
			LocationEdited: false,
			Accuracy:       1100,
		},
		Device: timeRegisterDeviceRequest{
			Manufacturer: "null",
			Model:        "null",
			Uuid: timeRegisterDeviceUuidRequest{
				Success:       signedUserSession.SignedUser.Success,
				Token:         signedUserSession.SignedUser.Token,
				ClientId:      signedUserSession.SignedUser.ClientId,
				Uuid:          signedUserSession.Uuid,
				Authorization: fmt.Sprintf("Bearer %v", signedUserSession.AuthorizationToken),
				Detail: timeRegisterDeviceUuidDetailRequest{
					Email:        signedUserSession.SignedUser.Detail.Login,
					SignInCount:  signedUserSession.SignedUser.Detail.SignInCount,
					LastSignInIp: signedUserSession.SignedUser.Detail.LastSignInIp,
					LastSignInAt: signedUserSession.SignedUser.Detail.LastSignInAt,
				},
			},
			Version: "null",
		},
		AppVersion: "0.10.32",
	}
}

func headersFrom(signedUserSession sign.SignedUserSession) map[string]string {
	return map[string]string{
		"api-version":  strconv.Itoa(session.API_VERSION),
		"access-token": signedUserSession.SignedUser.Token,
		"client":       signedUserSession.SignedUser.ClientId,
		"origin":       "https://app2.pontomais.com.br",
		"referer":      "https://app2.pontomais.com.br",
		"token":        signedUserSession.SignedUser.Token,
		"uid":          signedUserSession.SignedUser.Detail.Login,
		"uuid":         signedUserSession.Uuid,
		"user-agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.115 Safari/537.36",
	}
}
