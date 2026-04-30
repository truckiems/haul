package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempPlain(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "plain-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestEncryptFile_DecryptFile_RoundTrip(t *testing.T) {
	original := "DB_HOST=localhost\nDB_PASS=secret\n"
	src := writeTempPlain(t, original)
	dst := filepath.Join(t.TempDir(), "env.vault")
	passphrase := "test-passphrase"

	if err := EncryptFile(src, dst, passphrase); err != nil {
		t.Fatalf("EncryptFile() error: %v", err)
	}

	plaintext, err := DecryptFile(dst, passphrase)
	if err != nil {
		t.Fatalf("DecryptFile() error: %v", err)
	}

	if plaintext != original {
		t.Errorf("expected %q, got %q", original, plaintext)
	}
}

func TestDecryptFile_MissingHeader(t *testing.T) {
	src := writeTempPlain(t, "KEY=value\n")
	_, err := DecryptFile(src, "passphrase")
	if err == nil {
		t.Error("expected error for file without vault header")
	}
}

func TestDecryptFile_WrongPassphrase(t *testing.T) {
	src := writeTempPlain(t, "KEY=value\n")
	dst := filepath.Join(t.TempDir(), "env.vault")

	_ = EncryptFile(src, dst, "correct")

	_, err := DecryptFile(dst, "wrong")
	if err == nil {
		t.Error("expected error when decrypting with wrong passphrase")
	}
}

func TestIsVaultFile_True(t *testing.T) {
	src := writeTempPlain(t, "KEY=val\n")
	dst := filepath.Join(t.TempDir(), "env.vault")
	_ = EncryptFile(src, dst, "pass")

	if !IsVaultFile(dst) {
		t.Error("expected IsVaultFile to return true")
	}
}

func TestIsVaultFile_False(t *testing.T) {
	src := writeTempPlain(t, "KEY=val\n")
	if IsVaultFile(src) {
		t.Error("expected IsVaultFile to return false for plain file")
	}
}
