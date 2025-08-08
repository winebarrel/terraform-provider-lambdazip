package zip

import (
	arzip "archive/zip"
	"compress/flate"
	"io"
	"os"
	"sort"
)

func ZipFile(files []string, contents map[string]string, name string, level int) error {
	f, err := os.Create(name)

	if err != nil {
		return err
	}

	defer f.Close()

	return Zip(files, contents, f, level)
}

func Zip(files []string, contents map[string]string, out io.Writer, level int) error {
	w := arzip.NewWriter(out)

	w.RegisterCompressor(arzip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, level)
	})

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

	type content struct {
		name string
		data string
	}

	contentsList := []content{}

	for name, data := range contents {
		contentsList = append(contentsList, content{name: name, data: data})
	}

	sort.Slice(contentsList, func(i, j int) bool { return contentsList[i].name < contentsList[j].name })

	for _, c := range contentsList {
		f, err := w.Create(c.name)

		if err != nil {
			return err
		}

		_, err = f.Write([]byte(c.data))

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
