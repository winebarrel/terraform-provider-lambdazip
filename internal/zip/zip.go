package zip

import (
	arzip "archive/zip"
	"compress/flate"
	"io"
	"os"
	"sort"
	"strings"
)

func Strip(path string, n int) string {
	if n == 0 {
		return path
	}

	sep := string(os.PathSeparator)
	path = strings.TrimPrefix(path, sep)
	dirs := strings.Split(path, sep)

	if n >= len(dirs) {
		return ""
	}

	return strings.Join(dirs[n:], sep)
}

func ZipFile(files []string, contents map[string]string, name string, level int, strip int) error {
	f, err := os.Create(name)

	if err != nil {
		return err
	}

	defer f.Close()

	return Zip(files, contents, f, level, strip)
}

func Zip(files []string, contents map[string]string, out io.Writer, level int, strip int) error {
	w := arzip.NewWriter(out)

	w.RegisterCompressor(arzip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, level)
	})

	for _, name := range files {
		stripped := Strip(name, strip)

		if stripped == "" {
			continue
		}

		f, err := w.Create(stripped)

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
		stripped := Strip(c.name, strip)

		if stripped == "" {
			continue
		}

		f, err := w.Create(stripped)

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
