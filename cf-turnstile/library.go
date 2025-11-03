package cfturnstile

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	metapanic "meta/meta-panic"
	"net/http"
)

func IsTurnstileTokenValid(ctx context.Context, secret string, response string) (bool, error) {
	postTokenVerify := struct {
		Secret   string `json:"secret"`
		Response string `json:"response"`
	}{
		Secret:   secret,
		Response: response,
	}
	verifyUrl := "https://challenges.cloudflare.com/turnstile/v0/siteverify"
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
