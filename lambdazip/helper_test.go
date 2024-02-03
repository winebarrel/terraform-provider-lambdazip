package lambdazip_test

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"os"
)

func base64Sha256(buf []byte) string {
	sha256Sum := sha256.Sum256(buf)
	b64 := base64.StdEncoding.EncodeToString(sha256Sum[:])
	return b64
}

func isFileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func listZip(src []byte) ([]string, error) {
	buf := bytes.NewReader(src)
	r, err := zip.NewReader(buf, int64(len(src)))

	if err != nil {
		return nil, err
	}

	list := []string{}

	for _, file := range r.File {
		list = append(list, file.Name)
	}

	return list, nil
}
