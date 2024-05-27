package util

import (
	"testing"
	"time"
)

func TestNewToken(t *testing.T) {
    id := int64(123)
    secret := "secret"
    tokenTTL := time.Minute

    token, err := NewToken(id, secret, tokenTTL)
    if err != nil {
        t.Errorf("Error generating token: %v", err)
    }

    if token == "" {
        t.Errorf("Token is empty")
    }
}