package glob_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/terraform-provider-lambdazip/internal/glob"
)

func TestGlob(_t *testing.T) {
	assert := assert.New(_t)
	require := require.New(_t)

	cwd, _ := os.Getwd()
	os.Chdir(_t.TempDir())
	defer os.Chdir(cwd)

	os.Mkdir("app", 0755)
	os.Mkdir("app/lib", 0755)
	os.WriteFile("app/hello.rb", []byte("puts 'world'"), 0755)
	os.WriteFile("app/world.rb", []byte("puts 'hello'"), 0755)
	os.WriteFile("app/README.md", []byte("# hello.rb"), 0644)
	os.WriteFile("app/lib/const.rb", []byte("A = 100"), 0644)
	os.WriteFile("app/.gitignore", []byte("*.dylib"), 0644)

	tt := []struct {
		pattern  []string
		excludes []string
		expected []string
	}{
		{
			pattern:  []string{"**"},
			excludes: []string{},
			expected: []string{
				"app/.gitignore",
				"app/README.md",
				"app/hello.rb",
				"app/lib/const.rb",
				"app/world.rb",
			},
		},
		{
			pattern:  []string{"**/*.rb"},
			excludes: []string{},
			expected: []string{
				"app/hello.rb",
				"app/lib/const.rb",
				"app/world.rb",
			},
		},
		{
			pattern: []string{"**"},
			excludes: []string{
				"app/.*",
				"app/*.md",
			},
			expected: []string{
				"app/hello.rb",
				"app/lib/const.rb",
				"app/world.rb",
			},
		},
		{
			pattern:  []string{"app/*.rb", "**/*.md"},
			excludes: []string{"app/world.*"},
			expected: []string{
				"app/README.md",
				"app/hello.rb",
			},
		},
	}

	for _, t := range tt {
		files, err := glob.Glob(t.pattern, t.excludes)
		require.NoError(err)
		assert.Equal(t.expected, files)
	}
}
