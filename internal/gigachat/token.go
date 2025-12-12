package gigachat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type TokenProvider struct {
	token     string
	expiresAt time.Time
	mutex     sync.Mutex
	client    *http.Client
}

func NewTokenProvider() (*TokenProvider, error) {
	tlsCfg, err := LoadGigaChatTLS()
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsCfg,
		},
	}

	p := &TokenProvider{
		client: client,
	}

	if err := p.refresh(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *TokenProvider) refresh() error {
	authURL := os.Getenv("GIGACHAT_AUTH_URL")
	authKEY := os.Getenv("GIGACHAT_AUTH_KEY")

	p.mutex.Lock()
	defer p.mutex.Unlock()

	body := []byte("scope=GIGACHAT_API_PERS")

	req, err := http.NewRequest("POST", authURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+authKEY)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("RqUID", uuid.New().String())

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка запроса токена: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", slog.Any("error", err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("сервер вернул %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var t tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return err
	}

	p.token = t.AccessToken
	p.expiresAt = time.Now().Add(time.Duration(t.ExpiresIn) * time.Second)

	return nil
}

func (p *TokenProvider) GetToken() (string, error) {
	p.mutex.Lock()
	needsRefresh := time.Now().After(p.expiresAt.Add(-1 * time.Minute))
	p.mutex.Unlock()

	if needsRefresh {
		if err := p.refresh(); err != nil {
			return "", err
		}
	}

	return p.token, nil
}
