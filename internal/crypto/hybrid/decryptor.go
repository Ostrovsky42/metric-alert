// Package hybrid предоставляет функциональность для расшифровки данных с использованием гибридного шифрования.
package hybrid

import (
	"encoding/base64"

	"metric-alert/internal/crypto/asymmetric"
	"metric-alert/internal/crypto/symmetric"
	"metric-alert/internal/server/logger"
)

// Decryptor представляет декриптор, который может расшифровывать данные с
// использованием гибридного шифрования.
type Decryptor struct {
	Asymmetric asymmetric.Decryptor
	symmetric.Symmetric
	isIncluded bool
}

// NewDecryptor создает новый экземпляр Decryptor с предоставленным путем к асимметричному ключу для расшифровки.
// Если путь пуст, возвращается Decryptor c isIncluded false.
func NewDecryptor(path string) *Decryptor {
	if len(path) == 0 {
		return &Decryptor{}
	}

	decryptor, err := asymmetric.NewDecryptor(path)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("err create decryptor")
	}

	return &Decryptor{
		Asymmetric: decryptor,
		isIncluded: true,
	}
}

// DecryptData расшифровывает данные с использованием инициализационного вектора (initVector).
// Возвращает расшифрованные данные и любую ошибку, возникшую во время расшифровки.
func (d *Decryptor) DecryptData(cipherData []byte, initVector string) ([]byte, error) {
	iv, err := base64.StdEncoding.DecodeString(initVector)
	if err != nil {
		return nil, err
	}

	data, err := d.Symmetric.Decript(cipherData, iv)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// SetSymmetric устанавливает симметричный ключ для расшифровки с использованием предоставленного symmetricKey.
// Он расшифровывает симметричный ключ с использованием асимметричного ключа и инициализирует Decryptor с симметричным ключом.
func (d *Decryptor) SetSymmetric(symmetricKey string) error {
	encryptedKey, err := base64.StdEncoding.DecodeString(symmetricKey)
	if err != nil {
		return err
	}

	key, err := d.Asymmetric.Decrypt(encryptedKey)
	if err != nil {
		return err
	}
	s, err := symmetric.NewGCM(key)
	if err != nil {
		return err
	}

	d.Symmetric = s

	return nil
}

// IsNotIncluded возвращает true, если Decryptor не включен (например, если не был передан путь к ассиметричному ключу).
func (d *Decryptor) IsNotIncluded() bool {
	return !d.isIncluded
}
