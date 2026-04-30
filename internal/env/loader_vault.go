package env

import "fmt"

// LoadWithPassphrase loads environment variables from a file.
// If the file is a haul vault, it decrypts it first using the given passphrase.
// If passphrase is empty and the file is a vault, an error is returned.
func LoadWithPassphrase(path, passphrase string) (map[string]string, error) {
	if IsVaultFile(path) {
		if passphrase == "" {
			return nil, fmt.Errorf("%s is an encrypted vault but no passphrase was provided", path)
		}

		plaintext, err := DecryptFile(path, passphrase)
		if err != nil {
			return nil, fmt.Errorf("decrypting vault %s: %w", path, err)
		}

		return parseEnvString(plaintext)
	}

	return Load(path)
}

// parseEnvString parses env key=value pairs from an in-memory string.
func parseEnvString(content string) (map[string]string, error) {
	tf, cleanup, err := writeTempString(content)
	if err != nil {
		return nil, err
	}
	defer cleanup()
	return Load(tf)
}

// writeTempString writes content to a temp file and returns its path plus a cleanup func.
func writeTempString(content string) (string, func(), error) {
	import_os_once.Do(func() {})
	f, err := tempFileCreate()
	if err != nil {
		return "", nil, err
	}
	if _, err := f.WriteString(content); err != nil {
		f.Close()
		return "", nil, err
	}
	f.Close()
	return f.Name(), func() { removeFile(f.Name()) }, nil
}
