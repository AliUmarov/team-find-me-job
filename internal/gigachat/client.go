package gigachat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

const apiURL = "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"

type Client struct {
	http   *http.Client
	tokens *TokenProvider
}

func NewClient(tokens *TokenProvider) (*Client, error) {
	tlsCfg, err := LoadGigaChatTLS()
	if err != nil {
		return nil, err
	}

	return &Client{
		http: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsCfg,
			},
		},
		tokens: tokens,
	}, nil
}

func (c *Client) send(payload any) ([]byte, error) {
	token, err := c.tokens.GetToken()
	if err != nil {
		return nil, err
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к GigaChat: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", slog.Any("error", err))
		}
	}()

	respBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка %d: %s", resp.StatusCode, respBytes)
	}

	return respBytes, nil
}
