package main

import (
	"fmt"
	"log"

	"github.com/evzubkov/go-lknpd"
)

func main() {

	var phone string

	fmt.Println("Enter you phone (forma: 79XXXXXXXXX):")
	fmt.Scanf("%s\n", &phone)

	device := lknpd.NewDevice("")

	login, err := lknpd.LoginByPhone(lknpd.LoginRequest{Phone: phone, RequireTpToBeActive: true})
	if err != nil {
		log.Panic(err)
	}

	var code string

	fmt.Println("Enter verify code:")
	fmt.Scanf("%s\n", &code)

	verify, err := lknpd.VerifyCode(lknpd.VerifyCodeRequest{Phone: phone, Code: code,
		ChallengeToken: login.ChallengeToken, DeviceInfo: device})
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("DeviceId: %s\n", device.SourceDeviceId)
	fmt.Printf("Refresh token: %s\n", verify.RefreshToken)
}
