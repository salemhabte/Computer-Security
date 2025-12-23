package infrastructure

import (
	domain "security/domain"
	"security/config"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"os"
)

// AuditLogger writes JSON events to disk and encrypts with AES-GCM using LOG_ENC_KEY.
type AuditLogger struct {
	writer io.Writer
	aead   cipher.AEAD
}

func NewAuditLogger(path string) (domain.IAuditLogger, error) {
	key := []byte(config.LOG_ENC_KEY)
	if len(key) != 32 {
		return nil, errors.New("LOG_ENC_KEY must be 32 bytes for AES-256")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	return &AuditLogger{writer: f, aead: aead}, nil
}

func (l *AuditLogger) Log(event domain.AuditEvent) error {
	raw, err := json.Marshal(event)
	if err != nil {
		return err
	}
	nonce := make([]byte, l.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return err
	}
	ciphertext := l.aead.Seal(nonce, nonce, raw, nil)
	_, err = l.writer.Write(append(ciphertext, '\n'))
	return err
}

