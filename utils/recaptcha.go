package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	Score       float64  `json:"score,omitempty"`
	Action      string   `json:"action,omitempty"`
	ChallengeTS string   `json:"challenge_ts,omitempty"`
	HostName    string   `json:"hostname,omitempty"`
	ErrorCodes  []string `json:"error-codes,omitempty"`
}

func VerifyRecaptcha(token string) error {
	secret := os.Getenv("RECAPTCHA_SECRET")
	if secret == "" {
		return errors.New("recaptcha secret key not set")
	}

	reqBody := []byte("secret=" + secret + "&response" + token)
	resp, err := http.Post("http://www.google.com/recaptcha/api/siteverify", "application/x-www-form-urlencoded", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var results RecaptchaResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return err
	}

	if !results.Success {
		return errors.New("reCAPTCHA verification failed")
	}
	return nil
}
