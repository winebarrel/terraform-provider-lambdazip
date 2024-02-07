package glob

import (
	"sort"

	"github.com/bmatcuk/doublestar/v4"
)

func Glob(patterns []string, excludes []string, opts ...doublestar.GlobOption) ([]string, error) {
	fileSet := map[string]struct{}{}

	opts = append(opts,
		doublestar.WithFilesOnly(),
		doublestar.WithFailOnIOErrors(),
	)

	for _, pat := range patterns {
		matches, err := doublestar.FilepathGlob(pat, opts...)

		if err != nil {
			return nil, err
		}

		for _, f := range matches {
			fileSet[f] = struct{}{}
		}
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
