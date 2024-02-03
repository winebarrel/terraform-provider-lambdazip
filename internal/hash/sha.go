package hash

import (
	"crypto/sha256"
	"encoding/base64"
	"os"
)

func Base64Sha256(file string) (string, error) {
	buf, err := os.ReadFile(file)

	if err != nil {
		return "", err
	}

	sha256Sum := sha256.Sum256(buf)
	b64 := base64.StdEncoding.EncodeToString(sha256Sum[:])

	return b64, nil
}
