package hash_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/terraform-provider-lambdazip/internal/hash"
)

func TestBase64Md5(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	os.WriteFile("hello.rb", []byte("puts 'world'"), 0755)

	base64md5, err := hash.Base64Md5("hello.rb")
	require.NoError(err)
	assert.Equal("Ikxb6Lu+E5Jfz8LsrR3Aaw==", base64md5)
}
