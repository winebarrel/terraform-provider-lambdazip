package hash

import (
	"crypto/md5"
	"encoding/base64"
	"os"
)

func Base64Md5(file string) (string, error) {
	buf, err := os.ReadFile(file)

	if err != nil {
		return "", err
	}

	md5Sum := md5.Sum(buf)
	b64 := base64.StdEncoding.EncodeToString(md5Sum[:])

	return b64, nil
}
