package env

import (
	"strings"
	"testing"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	plaintext := "SECRET=hunter2\nAPI_KEY=abc123"
	passphrase := "my-secure-passphrase"

	encrypted, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	if encrypted == plaintext {
		t.Fatal("Encrypt() returned plaintext unchanged")
	}

	decrypted, err := Decrypt(encrypted, passphrase)
	if err != nil {
		t.Fatalf("Decrypt() error: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestEncrypt_DifferentCiphertexts(t *testing.T) {
	plaintext := "KEY=value"
	passphrase := "passphrase"

	a, _ := Encrypt(plaintext, passphrase)
	b, _ := Encrypt(plaintext, passphrase)

	if a == b {
		t.Error("expected different ciphertexts for same plaintext (random nonce)")
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	encrypted, _ := Encrypt("VALUE=secret", "correct")

	_, err := Decrypt(encrypted, "wrong")
	if err == nil {
		t.Error("expected error when decrypting with wrong passphrase")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	_, err := Decrypt("not-valid-base64!!!", "passphrase")
	if err == nil {
		t.Error("expected error for invalid base64 input")
	}
}

func TestDecrypt_TooShort(t *testing.T) {
	// A valid base64 string that decodes to fewer bytes than the nonce size.
	short := "aGk="olean
	_ = short
	tiny := strings.Repeat("A", 4)
	_, err := Decrypt(tiny, "passphrase")
	if err == nil {
		t.Error("expected error for ciphertext that is too short")
	}
}
