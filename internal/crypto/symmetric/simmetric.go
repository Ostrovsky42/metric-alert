package symmetric

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"io"
)

const (
	HeaderSymmetricKey         = "Symmetric-Key"
	HeaderInitializationVector = "Initialization-Vector"
	InitializationVectorSize   = 12
)

type Symmetric struct {
	cipher.AEAD
	key []byte
}

func GetRandKey() ([]byte, error) {
	symKey := make([]byte, aes.BlockSize)
	if _, err := rand.Read(symKey); err != nil {
		return nil, err
	}
	return symKey, nil
}

func NewGCM(key []byte) (Symmetric, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return Symmetric{}, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return Symmetric{}, err
	}

	return Symmetric{AEAD: aesGCM, key: key}, nil
}

func (s Symmetric) Encrypt(plainText []byte) ([]byte, []byte, error) {
	iv := make([]byte, InitializationVectorSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}
	cipherText := s.AEAD.Seal(nil, iv, plainText, nil)

	return cipherText, iv, nil
}

func (s Symmetric) Decript(cipherText, iv []byte) ([]byte, error) {
	return s.AEAD.Open(nil, iv, cipherText, nil)
}

func (s Symmetric) IsNotCreated() bool {
	return len(s.key) == 0
}
