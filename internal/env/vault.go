package env

import (
	"fmt"
	"os"
	"strings"
)

const vaultHeader = "# haul:encrypted\n"

// EncryptFile reads an env file, encrypts its contents, and writes a vault file.
func EncryptFile(src, dst, passphrase string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading source file: %w", err)
	}

	encrypted, err := Encrypt(string(data), passphrase)
	if err != nil {
		return fmt.Errorf("encrypting: %w", err)
	}

	contents := vaultHeader + encrypted + "\n"
	if err := os.WriteFile(dst, []byte(contents), 0600); err != nil {
		return fmt.Errorf("writing vault file: %w", err)
	}
	return nil
}

// DecryptFile reads a vault file and returns the decrypted env content.
func DecryptFile(src, passphrase string) (string, error) {
	data, err := os.ReadFile(src)
	if err != nil {
		return "", fmt.Errorf("reading vault file: %w", err)
	}

	content := string(data)
	if !strings.HasPrefix(content, vaultHeader) {
		return "", fmt.Errorf("file does not appear to be a haul vault (missing header)")
	}

	encrypted := strings.TrimSpace(strings.TrimPrefix(content, vaultHeader))
	plaintext, err := Decrypt(encrypted, passphrase)
	if err != nil {
		return "", fmt.Errorf("decrypting: %w", err)
	}

	return plaintext, nil
}

// IsVaultFile reports whether the given file path looks like an encrypted vault.
func IsVaultFile(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.HasPrefix(string(data), vaultHeader)
}
