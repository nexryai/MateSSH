package hostkey

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	gossh "golang.org/x/crypto/ssh"
)

type Keyring struct {
	//RSA鍵は大きいのでポインタ
	hostKeyRSA        *rsa.PrivateKey
	HostKeyRSAPem     string `json:"rsa"`
	hostKeyEd25519    ed25519.PrivateKey
	HostKeyEd25519Pem string          `json:"ed25519"`
	Signers           *[]gossh.Signer `json:"-"`
}

func (h *Keyring) Generate() error {
	// Generate host keys
	rsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	_, ed25519Key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	h.hostKeyRSA = rsaKey
	h.hostKeyEd25519 = ed25519Key

	// Generate PEM (RSA)
	h.HostKeyRSAPem, err = rsaPrivateKeyToPem(rsaKey)
	if err != nil {
		return err
	}

	// Generate PEM (Ed25519)
	h.HostKeyEd25519Pem, err = ed25519PrivateKeyToPem(ed25519Key)
	if err != nil {
		return err
	}

	// Generate signers
	err = h.GenSigners()
	if err != nil {
		return err
	}

	return nil
}

func (h *Keyring) GenSigners() error {
	rsaSigner, err := gossh.NewSignerFromKey(h.hostKeyRSA)
	if err != nil {
		return err
	}

	ed25519Signer, err := gossh.NewSignerFromKey(h.hostKeyEd25519)
	if err != nil {
		return err
	}

	h.Signers = &[]gossh.Signer{rsaSigner, ed25519Signer}
	return nil
}
