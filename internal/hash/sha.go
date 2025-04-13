package hash

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
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

func Sha256Map(files []string) (map[string]string, error) {
	m := map[string]string{}

	for _, f := range files {
		buf, err := os.ReadFile(f)

		if err != nil {
			return nil, err
		}

		sha256Sum := sha256.Sum256(buf)
		h := hex.EncodeToString(sha256Sum[:])
		m[f] = h
	}

	return m, nil
}

func ContentsSha256Map(contents map[string]string) map[string]string {
	m := map[string]string{}

	for name, data := range contents {
		sha256Sum := sha256.Sum256([]byte(data))
		h := hex.EncodeToString(sha256Sum[:])
		m[name] = h
	}

	return m
}
