package glob

import (
	"sort"

	"github.com/bmatcuk/doublestar/v4"
)

func Glob(pat string, excludes []string) ([]string, error) {
	matches, err := doublestar.FilepathGlob(pat,
		doublestar.WithFilesOnly(),
		doublestar.WithFailOnIOErrors())

	if err != nil {
		return nil, err
	}

	fileSet := map[string]struct{}{}

	for _, f := range matches {
		fileSet[f] = struct{}{}
	}

	for _, exPat := range excludes {
		exMatches, err := doublestar.FilepathGlob(exPat,
			doublestar.WithFilesOnly(),
			doublestar.WithFailOnIOErrors())

		if err != nil {
			return nil, err
		}

		for _, f := range exMatches {
			delete(fileSet, f)
		}
	}

	files := []string{}

	for f := range fileSet {
		files = append(files, f)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})

	return files, nil
}
