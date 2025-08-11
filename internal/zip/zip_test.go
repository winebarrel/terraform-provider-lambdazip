package zip_test

import (
	"bytes"
	"compress/flate"
	"os"
	"slices"
	"testing"

	arzip "archive/zip"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/terraform-provider-lambdazip/internal/zip"
)

func listZip(t *testing.T, src []byte) []string {
	t.Helper()
	require := require.New(t)

	buf := bytes.NewReader(src)
	r, err := arzip.NewReader(buf, int64(len(src)))
	require.NoError(err)
	list := []string{}

	for _, file := range r.File {
		list = append(list, file.Name)
	}

	slices.Sort(list)

	return list
}

func TestZip(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	os.WriteFile("hello.rb", []byte("puts 'world'"), 0755)
	os.WriteFile("world.rb", []byte("puts 'hello'"), 0755)

	var out bytes.Buffer
	err := zip.Zip([]string{"hello.rb", "world.rb"}, nil, &out, -1, 0)
	require.NoError(err)

	list := listZip(t, out.Bytes())
	assert.Equal([]string{"hello.rb", "world.rb"}, list)
}

func TestZipWithBestCompression(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	os.WriteFile("hello.rb", []byte("puts 'world'"), 0755)
	os.WriteFile("world.rb", []byte("puts 'hello'"), 0755)

	var out bytes.Buffer
	err := zip.Zip([]string{"hello.rb", "world.rb"}, nil, &out, flate.BestCompression, 0)
	require.NoError(err)

	list := listZip(t, out.Bytes())
	assert.Equal([]string{"hello.rb", "world.rb"}, list)
}

func TestZipContents(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	os.WriteFile("hello.rb", []byte("puts 'world'"), 0755)
	os.WriteFile("world.rb", []byte("puts 'hello'"), 0755)

	contents := map[string]string{
		"hello2.rb": "puts 'world2'",
		"world2.rb": "puts 'hello2'",
	}

	var out bytes.Buffer
	err := zip.Zip([]string{"hello.rb", "world.rb"}, contents, &out, -1, 0)
	require.NoError(err)

	list := listZip(t, out.Bytes())
	assert.Equal([]string{"hello.rb", "hello2.rb", "world.rb", "world2.rb"}, list)
}

func TestZipFile(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	os.WriteFile("hello.rb", []byte("puts 'world'"), 0755)
	os.WriteFile("world.rb", []byte("puts 'hello'"), 0755)

	err := zip.ZipFile([]string{"hello.rb", "world.rb"}, nil, "app.zip", -1, 0)
	require.NoError(err)
	buf, err := os.ReadFile("app.zip")
	require.NoError(err)

	list := listZip(t, buf)
	assert.Equal([]string{"hello.rb", "world.rb"}, list)
}

func TestZipFileWithContents(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	os.WriteFile("hello.rb", []byte("puts 'world'"), 0755)
	os.WriteFile("world.rb", []byte("puts 'hello'"), 0755)

	contents := map[string]string{
		"hello2.rb": "puts 'world2'",
		"world2.rb": "puts 'hello2'",
	}

	err := zip.ZipFile([]string{"hello.rb", "world.rb"}, contents, "app.zip", -1, 0)
	require.NoError(err)
	buf, err := os.ReadFile("app.zip")
	require.NoError(err)

	list := listZip(t, buf)
	assert.Equal([]string{"hello.rb", "hello2.rb", "world.rb", "world2.rb"}, list)
}

func TestZipWithRemovePrefix(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	os.Mkdir("app", 0755)
	os.WriteFile("app/hello.rb", []byte("puts 'world'"), 0755)
	os.WriteFile("app/world.rb", []byte("puts 'hello'"), 0755)

	var out bytes.Buffer
	err := zip.Zip([]string{"app/hello.rb", "app/world.rb"}, nil, &out, -1, 1)
	require.NoError(err)

	list := listZip(t, out.Bytes())
	assert.Equal([]string{"hello.rb", "world.rb"}, list)
}

func TestZipContentsWithRemovePrefix(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	os.Mkdir("app", 0755)
	os.WriteFile("app/hello.rb", []byte("puts 'world'"), 0755)
	os.WriteFile("app/world.rb", []byte("puts 'hello'"), 0755)

	contents := map[string]string{
		"app/hello2.rb": "puts 'world2'",
		"app/world2.rb": "puts 'hello2'",
	}

	var out bytes.Buffer
	err := zip.Zip([]string{"app/hello.rb", "app/world.rb"}, contents, &out, -1, 1)
	require.NoError(err)

	list := listZip(t, out.Bytes())
	assert.Equal([]string{"hello.rb", "hello2.rb", "world.rb", "world2.rb"}, list)
}

func TestStrip(t *testing.T) {
	assert := assert.New(t)

	tt := []struct {
		path     string
		n        int
		expected string
	}{
		{path: "foo/bar/zoo", n: 0, expected: "foo/bar/zoo"},
		{path: "foo/bar/zoo", n: 1, expected: "bar/zoo"},
		{path: "foo/bar/zoo", n: 2, expected: "zoo"},
		{path: "foo/bar/zoo", n: 3, expected: ""},
		{path: "foo/bar/zoo", n: 4, expected: ""},
		{path: "/foo/bar/zoo", n: 0, expected: "/foo/bar/zoo"},
		{path: "/foo/bar/zoo", n: 1, expected: "bar/zoo"},
		{path: "/foo/bar/zoo", n: 2, expected: "zoo"},
		{path: "/foo/bar/zoo", n: 3, expected: ""},
		{path: "/foo/bar/zoo", n: 4, expected: ""},
	}

	for _, t := range tt {
		actual := zip.Strip(t.path, t.n)
		assert.Equal(t.expected, actual)
	}
}
