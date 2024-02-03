package cmd_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/terraform-provider-lambdazip/internal/cmd"
)

func TestRun_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	os.WriteFile("hello.rb", []byte("puts 'world'"), 0755)
	os.WriteFile("world.rb", []byte("puts 'hello'"), 0755)

	out, err := cmd.Run("ls ./")
	require.NoError(err)
	assert.Equal(`hello.rb
world.rb`, out)
}

func TestRun_Err(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	os.WriteFile("hello.rb", []byte("puts 'world'"), 0755)
	os.WriteFile("world.rb", []byte("puts 'hello'"), 0755)

	out, err := cmd.Run("ls /not/exist")
	require.Error(err)
	assert.NotEmpty(out)
}
