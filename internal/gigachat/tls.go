package gigachat

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func LoadGigaChatTLS() (*tls.Config, error) {
	caCert, err := os.ReadFile("certs/gigachat-bundle.pem")
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения gigachat-bundle.pem: %w", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("не удалось загрузить сертификаты в trust pool")
	}

	return &tls.Config{
		RootCAs:    caPool,
		MinVersion: tls.VersionTLS12,
	}, nil
}
