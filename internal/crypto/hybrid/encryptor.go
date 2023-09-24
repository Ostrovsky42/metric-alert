// Package hybrid предоставляет функциональность для шифрования данных с использованием гибридного шифрования.
package hybrid

import (
	"encoding/base64"
	"metric-alert/internal/crypto/asymmetric"
	"metric-alert/internal/crypto/symmetric"
)

// Encryptor представляет шифратор, который может шифровать данные с использованием гибридного шифрования.
type Encryptor struct {
	Asymmetric   *asymmetric.Encryptor
	Symmetric    *symmetric.Symmetric
	EncryptedKey string
	isIncluded   bool
}

// NewEncryptor создает новый экземпляр Encryptor с предоставленным путем к асимметричному ключу для шифрования.
// Если путь пуст, возвращается Encryptor c isIncluded false.
// Ключи шифрования генерируются случайным образом.
func NewEncryptor(path string) (*Encryptor, error) {
	if len(path) == 0 {
		return &Encryptor{}, nil
	}

	encryptor, err := asymmetric.NewEncryptor(path)
	if err != nil {
		return nil, err
	}

	key, err := symmetric.GetRandKey()
	if err != nil {
		return nil, err
	}

	encKey, err := encryptor.Encrypt(key)
	if err != nil {
		return nil, err
	}

	s, err := symmetric.NewGCM(key)
	if err != nil {
		return nil, err
	}

	return &Encryptor{
		Asymmetric:   encryptor,
		Symmetric:    s,
		EncryptedKey: base64.StdEncoding.EncodeToString(encKey),
		isIncluded:   true,
	}, nil
}

// Included возвращает false, если Encryptor выключен (например, если не был передан путь к ассиметричному ключу).
func (e *Encryptor) Included() bool {
	return e.isIncluded
}
