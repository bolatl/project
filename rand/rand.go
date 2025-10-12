package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const SessionTokenBytes = 32

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("error in Bytes func: %w", err)
	}
	if nRead < n {
		return nil, fmt.Errorf("not enough bytes generated in Bytes func")
	}
	return b, nil
}

func String(n int) (string, error) {
	s, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("error in String func: %w", err)
	}
	return base64.URLEncoding.EncodeToString(s), err
}

func SessionToken() (string, error) {
	return String(SessionTokenBytes)
}
