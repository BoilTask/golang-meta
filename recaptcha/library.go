package recaptcha

import (
	"context"
	"encoding/json"
	"io"
	metapanic "meta/meta-panic"
	"net/http"
	"net/url"
)

func IsRecaptchaTokenValid(ctx context.Context, secret string, response string, ip string) (bool, error) {
	verifyUrl := "https://recaptcha.net/recaptcha/api/siteverify"
	formData := url.Values{}
	formData.Set("secret", secret)
	formData.Set("response", response)
	// formData.Set("remoteip", ip)

	resp, err := http.PostForm(verifyUrl, formData)
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
		Success     bool     `json:"success"`
		ChallengeTS string   `json:"challenge_ts"`
		Hostname    string   `json:"hostname"`
		ErrorCodes  []string `json:"error-codes"`
	}
	err = json.NewDecoder(resp.Body).Decode(&verifyResponseData)
	if err != nil {
		return false, err
	}
	return verifyResponseData.Success, nil
}
