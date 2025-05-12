package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"
)

type Encrypt interface {
	Encrypt(value string) (string, error)
	GenerateBytes(n int) ([]byte, error)
	GenerateString(n int) (string, error)
	GenerateHash(value string) (string, error)
}

type DefaultEncrypt struct{}

const encryptKey = "!S@perS!kr3TKe!eY"

func (e *DefaultEncrypt) Encrypt(value string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(value + encryptKey))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (e *DefaultEncrypt) GenerateBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (e *DefaultEncrypt) GenerateString(n int) (string, error) {
	b, err := e.GenerateBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}

func (e *DefaultEncrypt) GenerateHash(value string) (string, error) {
	str, err := e.GenerateString(16)

	if err != nil {
		return "", err
	}

	hash := sha256.New()
	_, err = hash.Write([]byte(time.Now().String() + str + value + encryptKey))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
