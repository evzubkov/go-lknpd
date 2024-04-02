package lknpd

import (
	"log"
	"time"

	"github.com/denisbrodbeck/machineid"
)

const (
	defaultAgent      = "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36"
	defaultSourceType = "WEB"
	defaultAppVersion = "1.0.0"
)

type Client struct {
	timezone            string
	refreshToken        string
	accessToken         string
	accessTokenExpireIn time.Time
	device              Device
}

func NewClient(timezone, deviceId, refreshToken string) *Client {
	return &Client{timezone: timezone, refreshToken: refreshToken, device: *NewDevice(deviceId)}
}

func (o *Client) CheckTokenExpireIn() {
	if time.Now().After(o.accessTokenExpireIn) {
		o.RefreshToken()
	}
}

func getDeviceId() string {
	id, err := machineid.ProtectedID("my-tax")
	if err != nil {
		log.Panic(err)
	}

	if len(id) > 21 {
		return id[:21]
	}

	return id
}

type Device struct {
	SourceDeviceId string `json:"sourceDeviceId"`
	SourceType     string `json:"sourceType"`
	AppVersion     string `json:"appVersion"`
	MetaDetails    struct {
		UserAgent string `json:"userAgent"`
	} `json:"metaDetails"`
}

func NewDevice(deviceId string) *Device {

	if deviceId == "" {
		deviceId = getDeviceId()
	}

	return &Device{SourceDeviceId: deviceId, SourceType: defaultSourceType,
		AppVersion: defaultAppVersion, MetaDetails: struct {
			UserAgent string "json:\"userAgent\""
		}{UserAgent: defaultAgent}}
}
