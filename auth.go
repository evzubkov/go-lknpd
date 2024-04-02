package lknpd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type (
	LoginRequest struct {
		Phone               string `json:"phone"`
		RequireTpToBeActive bool   `json:"requireTpToBeActive"`
	}

	LoginResponse struct {
		ChallengeToken string    `json:"challengeToken"`
		ExpireDate     time.Time `json:"expireDate"`
		ExpireIn       int       `json:"expireIn"`
	}
)

func LoginByPhone(payload LoginRequest) (*LoginResponse, error) {

	client := resty.New()
	resp, err := client.R().SetBody(payload).
		SetHeader("Content-Type", "application/json").
		SetHeader("Referrer", "https://lknpd.nalog.ru/").
		SetHeader("Referrer-Policy", "strict-origin-when-cross-origin").
		Post("https://lknpd.nalog.ru/api/v2/auth/challenge/sms/start")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == 200 {
		var result LoginResponse
		err := json.Unmarshal(resp.Body(), &result)
		if err != nil {
			return nil, err
		}
		return &result, nil
	} else {
		err = fmt.Errorf("status code: %d. msg: %s", resp.StatusCode(), resp.RawBody())
	}

	return nil, err
}

type (
	VerifyCodeRequest struct {
		Phone          string      `json:"phone"`
		Code           string      ` json:"code"`
		ChallengeToken string      `json:"challengeToken"`
		DeviceInfo     interface{} `json:"deviceInfo"`
	}

	VerifyCodeResponse struct {
		RefreshToken string `json:"refreshToken"`
	}
)

func VerifyCode(payload VerifyCodeRequest) (*VerifyCodeResponse, error) {

	client := resty.New()
	resp, err := client.R().SetBody(payload).
		SetHeader("Content-Type", "application/json").
		SetHeader("Referrer", "https://lknpd.nalog.ru/").
		SetHeader("Referrer-Policy", "strict-origin-when-cross-origin").
		Post("https://lknpd.nalog.ru/api/v1/auth/challenge/sms/verify")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == 200 {
		var result VerifyCodeResponse
		err := json.Unmarshal(resp.Body(), &result)
		if err != nil {
			return nil, err
		}

		return &result, nil
	} else {
		err = fmt.Errorf("status code: %d. msg: %s", resp.StatusCode(), resp.RawBody())
	}

	return nil, err
}

type (
	RefreshTokenRequest struct {
		RefreshToken string      `json:"refreshToken"`
		DeviceInfo   interface{} `json:"deviceInfo"`
	}

	RefreshTokenResponese struct {
		RefreshToken  string    `json:"refreshToken"`
		Token         string    `json:"token"`
		TokenExpireIn time.Time `json:"tokenExpireIn"`
	}
)

func (o *Client) RefreshToken() (err error) {

	client := resty.New()
	resp, err := client.R().SetBody(RefreshTokenRequest{RefreshToken: o.refreshToken, DeviceInfo: o.device}).
		SetHeader("Content-Type", "application/json").
		SetHeader("Referrer", "https://lknpd.nalog.ru/").
		SetHeader("Referrer-Policy", "strict-origin-when-cross-origin").
		Post("https://lknpd.nalog.ru/api/v1/auth/token")
	if err != nil {
		return
	}

	if resp.StatusCode() == 200 {
		var result RefreshTokenResponese
		if err = json.Unmarshal(resp.Body(), &result); err != nil {
			return
		}
		o.accessToken = result.Token
		o.accessTokenExpireIn = result.TokenExpireIn
	} else {
		err = fmt.Errorf("status code: %d. msg: %s", resp.StatusCode(), resp.RawBody())
	}

	return
}
