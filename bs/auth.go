package bs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
)

type JWTToken struct {
	Token   string
	Refresh string
}

type LoginCredentials struct {
	Email    string
	Password string
}

func (bsRef *Bigscreen) login() {
	ctx := *bsRef
	reqBody := fmt.Sprintf(`{"email":"%s", "password": "%s"}`, ctx.Credentials.Email, ctx.Credentials.Password)

	_, resp := bsRef.request(
		ctx.HostAccounts+"/login",
		"POST",
		map[string]string{
			"content-type":            "application/json",
			"x-bigscreen-system-info": base64.StdEncoding.EncodeToString([]byte(ctx.DeviceInfo)),
		},
		reqBody,
	)

	(*bsRef).JWT.Token = resp.Header.Get("x-access-token")
	(*bsRef).JWT.Refresh = resp.Header.Get("x-refresh-token")
}

func (bsRef *Bigscreen) verify() {
	ctx := *bsRef

	if ctx.JWT.Refresh == "" {
		bsRef.login()
	}

	respBody, resp := bsRef.request(
		ctx.HostAccounts+"/verify",
		"GET",
		make(map[string]string),
		"",
	)
	_ = len(respBody)
	if resp.StatusCode == 401 {
		bsRef.renew(resp.Header.Get("x-bigscreen-nonce"))
		bsRef.verify()
	}

}

func (bsRef *Bigscreen) renew(nonce string) {
	ctx := *bsRef
	respBody, resp := bsRef.request(
		ctx.HostAccounts+"/renew",
		"GET",
		map[string]string{
			"x-bigscreen-system-info": base64.StdEncoding.EncodeToString([]byte(ctx.DeviceInfo)),
			"x-refresh-token":         ctx.JWT.Refresh,
			"x-bigscreen-nonce":       nonce,
		},
		"",
	)

	if resp.StatusCode == 401 {
		var apiErr ApiError
		err := json.Unmarshal(respBody, &apiErr)
		if err != nil {
			log.Panic(err.Error())
		}

		switch apiErr.Code {
		case 2:
			log.Println(apiErr.Message)
			bsRef.login()
			bsRef.renew(nonce)
		default:
			log.Panic(respBody)
		}
	}

	if resp.StatusCode == 200 {
		(*bsRef).JWT.Token = resp.Header.Get("x-access-token")
	}
}
