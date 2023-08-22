// Пакет hasher предоставляет функции для генерации хешей на основе алгоритма SHA-256.
package hasher

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashBuilder интерфейс генератора хешей.
type HashBuilder interface {
	GetHash(data []byte) string // GetHash генерирует хеш для переданных данных и возвращает его в виде строки.
	IsNotActive() bool          // IsNotActive возвращает true, если генератор хешей не активен (не имеет ключа).
}

// HashGenerator представляет генератор хешей с использованием заданного ключа.
type HashGenerator struct {
	signKey  []byte
	isActive bool
}

// NewHashGenerator создает и возвращает новый экземпляр HashGenerator с заданным ключом.
// Если ключ не задан (пустая строка), генератор считается неактивным.
func NewHashGenerator(key string) *HashGenerator {
	return &HashGenerator{
		isActive: len(key) > 0,
		signKey:  []byte(key),
	}
}

// GetHash генерирует хеш для переданных данных с использованием алгоритма SHA-256 и возвращает хеш в виде строки.
func (h *HashGenerator) GetHash(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	signedHash := hash.Sum(h.signKey)
	return hex.EncodeToString(signedHash)
}

// IsNotActive возвращает true, если генератор хешей не активен (не имеет ключа).
func (h *HashGenerator) IsNotActive() bool {
	return !h.isActive
}
