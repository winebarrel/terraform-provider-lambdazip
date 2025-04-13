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

func TestSha256Map(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	os.WriteFile("hello.rb", []byte("puts 'world'"), 0755)
	os.WriteFile("world.rb", []byte("puts 'hello'"), 0755)

	m, err := hash.Sha256Map([]string{"hello.rb", "world.rb"})
	require.NoError(err)

	assert.Equal(map[string]string{
		"hello.rb": "06db2c7a260efaf6e2e3f4c635c83506f1f40f6d3898e0e6025e3e55f44ddebe",
		"world.rb": "293c10e07909b3a823d7d2ba87c6cdf7400c9ed70132c2c952d7c8147d945a74",
	}, m)
}

func TestContentsSha256Map(t *testing.T) {
	assert := assert.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	os.WriteFile("hello.rb", []byte("puts 'world'"), 0755)
	os.WriteFile("world.rb", []byte("puts 'hello'"), 0755)

	m := hash.ContentsSha256Map(map[string]string{
		"hello.rb": "puts 'world'",
		"world.rb": "puts 'hello'",
	})

	assert.Equal(map[string]string{
		"hello.rb": "06db2c7a260efaf6e2e3f4c635c83506f1f40f6d3898e0e6025e3e55f44ddebe",
		"world.rb": "293c10e07909b3a823d7d2ba87c6cdf7400c9ed70132c2c952d7c8147d945a74",
	}, m)
}
