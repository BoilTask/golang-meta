package recaptcha

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	metapanic "meta/meta-panic"
	"net/http"
)

func IsRecaptchaTokenValid(ctx context.Context, secret string, response string, ip string) (bool, error) {
	postTokenVerify := struct {
		Secret   string `json:"secret"`
		Response string `json:"response"`
		RemoteIP string `json:"remoteip"`
	}{
		Secret:   secret,
		Response: response,
		RemoteIP: ip,
	}
	verifyUrl := "https://recaptcha.net/recaptcha/api/siteverify"
	jsonData, err := json.Marshal(postTokenVerify)
	if err != nil {
		return false, err
	}
	resp, err := http.Post(verifyUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			metapanic.ProcessError(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}
	var verifyResponseData struct {
		Success bool `json:"success"`
	}
	err = json.NewDecoder(resp.Body).Decode(&verifyResponseData)
	if err != nil {
		return false, err
	}
	return verifyResponseData.Success, nil
}
