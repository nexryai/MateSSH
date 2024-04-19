package hostkey

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"reflect"
	"testing"
)

func TestRSAKeyEncodeAndDecode(t *testing.T) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Fatal(err)
	}

	pemKey, err := rsaPrivateKeyToPem(rsaKey)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(fmt.Sprintf("RSA key: %s", pemKey))
	}

	decodedKey, err := rsaPemToPrivateKey(pemKey)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(rsaKey, decodedKey) {
		t.Fatal("Decoded key is not equal to the original key")
	}
}

func TestEd25519KeyEncodeAndDecode(t *testing.T) {
	_, ed25519Key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	pemKey, err := ed25519PrivateKeyToPem(ed25519Key)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(fmt.Sprintf("Ed25519 key: %s", pemKey))
	}

	decodedKey, err := ed25519PemToPrivateKey(pemKey)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(ed25519Key, decodedKey) {
		t.Fatal("Decoded key is not equal to the original key")
	}
}
