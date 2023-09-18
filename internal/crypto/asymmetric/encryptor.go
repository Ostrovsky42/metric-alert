package asymmetric

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

type Encryptor struct {
	*rsa.PublicKey
}

func NewEncryptor(path string) (Encryptor, error) {
	key, err := readPublicKeyFromFile(path)
	if err != nil {
		return Encryptor{}, err
	}

	return Encryptor{PublicKey: key}, nil
}

func (e *Encryptor) Encrypt(plaintext []byte) ([]byte, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, e.PublicKey, plaintext)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

func readPublicKeyFromFile(filename string) (*rsa.PublicKey, error) {
	publicKeyPEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicKeyPEM)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid PEM block")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return rsaPublicKey, nil
}
