package hostkey

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	gossh "golang.org/x/crypto/ssh"
)

type Keyring struct {
	//RSA鍵は大きいのでポインタ
	HostKeyRSA     *rsa.PrivateKey
	HostKeyEd25519 ed25519.PrivateKey
	Signers        *[]gossh.Signer
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

	h.HostKeyRSA = rsaKey
	h.HostKeyEd25519 = ed25519Key

	// Generate signers
	err = h.GenSigners()
	if err != nil {
		return err
	}

	return nil
}

func (h *Keyring) GenSigners() error {
	rsaSigner, err := gossh.NewSignerFromKey(h.HostKeyRSA)
	if err != nil {
		return err
	}

	ed25519Signer, err := gossh.NewSignerFromKey(h.HostKeyEd25519)
	if err != nil {
		return err
	}

	h.Signers = &[]gossh.Signer{rsaSigner, ed25519Signer}
	return nil
}
