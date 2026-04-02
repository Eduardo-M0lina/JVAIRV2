package sms

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TwilioSender implementa SMSSender usando la API REST de Twilio
type TwilioSender struct {
	sid        string
	authToken  string
	fromNumber string
	httpClient *http.Client
}

// twilioErrorResponse representa la respuesta de error de la API de Twilio
type twilioErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// NewTwilioSender crea un nuevo cliente Twilio usando credenciales de la tabla settings
func NewTwilioSender(sid, authToken, fromNumber string) *TwilioSender {
	return &TwilioSender{
		sid:        sid,
		authToken:  authToken,
		fromNumber: fromNumber,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// SendSMS envía un SMS via Twilio REST API
func (s *TwilioSender) SendSMS(ctx context.Context, to string, body string) error {
	apiURL := fmt.Sprintf(
		"https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json",
		s.sid,
	)

	data := url.Values{}
	data.Set("To", to)
	data.Set("From", s.fromNumber)
	data.Set("Body", body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("error al crear solicitud Twilio: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(s.sid + ":" + s.authToken))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error al contactar API de Twilio: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	responseBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var twilioErr twilioErrorResponse
		if json.Unmarshal(responseBody, &twilioErr) == nil && twilioErr.Message != "" {
			return fmt.Errorf("twilio error %d: %s (código: %d)", resp.StatusCode, twilioErr.Message, twilioErr.Code)
		}
		return fmt.Errorf("twilio respondió con estado %d: %s", resp.StatusCode, string(responseBody))
	}

	return nil
}
