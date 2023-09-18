package asymmetric

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

type Decryptor struct {
	*rsa.PrivateKey
}

func NewDecryptor(path string) (Decryptor, error) {
	key, err := readPrivateKeyFromFile(path)
	if err != nil {
		return Decryptor{}, err
	}

	return Decryptor{key}, nil
}

func (d *Decryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, d.PrivateKey, ciphertext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func readPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	privateKeyPEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("invalid PEM block")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}

	return rsaPrivateKey, nil
}
