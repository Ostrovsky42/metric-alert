package hasher

import (
	"crypto/sha256"
	"encoding/hex"
)

type HashBuilder interface {
	GetHash(data []byte) string
	IsNotActive() bool
}

type HashGenerator struct {
	signKey  []byte
	isActive bool
}

func NewHashGenerator(key string) *HashGenerator {
	return &HashGenerator{
		isActive: len(key) > 0,
		signKey:  []byte(key),
	}
}

func (h *HashGenerator) GetHash(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	signedHash := hash.Sum(h.signKey)
	return hex.EncodeToString(signedHash)
}

func (h *HashGenerator) IsNotActive() bool {
	return !h.isActive
}
