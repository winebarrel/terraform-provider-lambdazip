package zip

import (
	arzip "archive/zip"
	"io"
	"os"
)

func ZipFile(files []string, name string) error {
	f, err := os.Create(name)

	if err != nil {
		return err
	}

	defer f.Close()

	return Zip(files, f)
}

func Zip(files []string, out io.Writer) error {
	w := arzip.NewWriter(out)

	for _, name := range files {
		f, err := w.Create(name)

		if err != nil {
			return err
		}

		buf, err := os.ReadFile(name)

		if err != nil {
			return err
		}

		_, err = f.Write(buf)

		if err != nil {
			return err
		}
	}

	err := w.Close()

	if err != nil {
		return err
	}

	return nil
}
