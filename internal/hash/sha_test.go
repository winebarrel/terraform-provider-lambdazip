package hash_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/terraform-provider-lambdazip/internal/hash"
)

func TestBase64Sha256(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	os.WriteFile("hello.rb", []byte("puts 'world'"), 0755)

	base64Sha256, err := hash.Base64Sha256("hello.rb")
	require.NoError(err)
	assert.Equal("BtsseiYO+vbi4/TGNcg1BvH0D204mODmAl4+VfRN3r4=", base64Sha256)
}
