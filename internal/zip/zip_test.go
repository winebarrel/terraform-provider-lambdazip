package zip_test

import (
	"bytes"
	"os"
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
	err := zip.Zip([]string{"hello.rb", "world.rb"}, &out)
	require.NoError(err)

	list := listZip(t, out.Bytes())
	assert.Equal([]string{"hello.rb", "world.rb"}, list)
}

func TestZipFile(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	os.WriteFile("hello.rb", []byte("puts 'world'"), 0755)
	os.WriteFile("world.rb", []byte("puts 'hello'"), 0755)

	err := zip.ZipFile([]string{"hello.rb", "world.rb"}, "app.zip")
	require.NoError(err)
	buf, err := os.ReadFile("app.zip")
	require.NoError(err)

	list := listZip(t, buf)
	assert.Equal([]string{"hello.rb", "world.rb"}, list)
}
