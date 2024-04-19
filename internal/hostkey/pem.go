package hostkey

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func rsaPrivateKeyToPem(key *rsa.PrivateKey) (string, error) {
	var privateKeyBuffer bytes.Buffer

	err := pem.Encode(&privateKeyBuffer, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	if err != nil {
		return "", err
	}

	return privateKeyBuffer.String(), nil
}

func rsaPemToPrivateKey(pemKey string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func ed25519PrivateKeyToPem(key ed25519.PrivateKey) (string, error) {
	var privateKeyBuffer bytes.Buffer

	b, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", err
	}

	err = pem.Encode(&privateKeyBuffer, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: b,
	})
	if err != nil {
		return "", err
	}

	return privateKeyBuffer.String(), nil
}

func ed25519PemToPrivateKey(pemKey string) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey.(ed25519.PrivateKey), nil
}
