package lambdazip_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFiles_basic(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	os.Mkdir("app", 0755)
	os.Mkdir("app/lib", 0755)
	os.WriteFile("app/hello.rb", []byte("puts 'world'"), 0755)
	os.WriteFile("app/world.rb", []byte("puts 'hello'"), 0755)
	os.WriteFile("app/README.md", []byte("# hello.rb"), 0644)
	os.WriteFile("app/lib/const.rb", []byte("A = 100"), 0644)

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			// Step 1 =====================================================
			{
				Config: `
					resource "lambdazip_file" "my_app" {
						base_dir      = "app"
						sources       = ["**/*.rb"]
						excludes      = [".*", "README.md"]
						output        = "my-app.zip"
						before_create = "touch exec.txt"

						triggers = {
							hello_rb = filesha256("app/hello.rb"),
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base_dir", "app"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.#", "1"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.0", "**/*.rb"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.#", "2"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.0", ".*"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.1", "README.md"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "output", "my-app.zip"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "before_create", "touch exec.txt"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "triggers.hello_rb", "06db2c7a260efaf6e2e3f4c635c83506f1f40f6d3898e0e6025e3e55f44ddebe"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base64sha256", "a6p1iUHprlLD6/MDM7kHa9dhCuOmcPiBmU+JShrJ4Ro="),
					func(*terraform.State) error {
						buf, err := os.ReadFile("my-app.zip")
						require.NoError(err)
						assert.Equal("a6p1iUHprlLD6/MDM7kHa9dhCuOmcPiBmU+JShrJ4Ro=", base64Sha256(buf))
						assert.True(isFileExists("app/exec.txt"))
						list, err := listZip(buf)
						require.NoError(err)
						assert.Equal([]string{"hello.rb", "lib/const.rb", "world.rb"}, list)
						return nil
					},
				),
			},
			// Step 2 =====================================================
			{
				Config: `
					resource "lambdazip_file" "my_app" {
						base_dir      = "app"
						sources       = ["**/*.rb"]
						excludes      = [".*", "README.md"]
						output        = "my-app.zip"
						before_create = "touch exec.txt"

						triggers = {
							hello_rb = filesha256("app/hello.rb"),
						}
					}
				`,
				PreConfig: func() {
					err := os.Remove("my-app.zip")
					require.NoError(err)
					err = os.Remove("app/exec.txt")
					require.NoError(err)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base_dir", "app"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.#", "1"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.0", "**/*.rb"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.#", "2"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.0", ".*"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.1", "README.md"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "output", "my-app.zip"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "before_create", "touch exec.txt"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "triggers.hello_rb", "06db2c7a260efaf6e2e3f4c635c83506f1f40f6d3898e0e6025e3e55f44ddebe"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base64sha256", "a6p1iUHprlLD6/MDM7kHa9dhCuOmcPiBmU+JShrJ4Ro="),
					func(*terraform.State) error {
						assert.False(isFileExists("my-app.zip"))
						assert.False(isFileExists("app/exec.txt"))
						return nil
					},
				),
			},
			// Step 3 =====================================================
			{
				Config: `
					resource "lambdazip_file" "my_app" {
						base_dir      = "app"
						sources       = ["**/*.rb"]
						excludes      = [".*", "README.md"]
						output        = "my-app.zip"
						before_create = "touch exec.txt"

						triggers = {
							hello_rb = filesha256("app/hello.rb"),
						}
					}
				`,
				PreConfig: func() {
					err := os.WriteFile("app/hello.rb", []byte("print 'world'"), 0755)
					require.NoError(err)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base_dir", "app"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.#", "1"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.0", "**/*.rb"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.#", "2"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.0", ".*"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.1", "README.md"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "output", "my-app.zip"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "before_create", "touch exec.txt"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "triggers.hello_rb", "6740287d0049734d6fe501a11d8572ba1befdc690d08d891db539d2f8a9d7273"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base64sha256", "55024qLWGJdd9NUOXG8Y/4xpeWuckpVC9VJE1ZZPQtA="),
					func(*terraform.State) error {
						buf, err := os.ReadFile("my-app.zip")
						require.NoError(err)
						assert.Equal("55024qLWGJdd9NUOXG8Y/4xpeWuckpVC9VJE1ZZPQtA=", base64Sha256(buf))
						assert.True(isFileExists("app/exec.txt"))
						list, err := listZip(buf)
						require.NoError(err)
						assert.Equal([]string{"hello.rb", "lib/const.rb", "world.rb"}, list)
						return nil
					},
				),
			},
			// Step 4 =====================================================
			{
				Config: `
					resource "lambdazip_file" "my_app" {
						base_dir      = "app"
						sources       = ["**/*.rb"]
						excludes      = [".*", "README.md"]
						output        = "my-app.zip"
						before_create = "touch exec.txt"

						triggers = {
							hello_rb = filesha256("app/hello.rb"),
						}
					}
				`,
				PreConfig: func() {
					err := os.Remove("my-app.zip")
					require.NoError(err)
					err = os.Remove("app/exec.txt")
					require.NoError(err)
					err = os.WriteFile("app/world.rb", []byte("print 'hello'"), 0755)
					require.NoError(err)
					err = os.Remove("app/lib/const.rb")
					require.NoError(err)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base_dir", "app"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.#", "1"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.0", "**/*.rb"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.#", "2"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.0", ".*"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.1", "README.md"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "output", "my-app.zip"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "before_create", "touch exec.txt"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "triggers.hello_rb", "6740287d0049734d6fe501a11d8572ba1befdc690d08d891db539d2f8a9d7273"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base64sha256", "55024qLWGJdd9NUOXG8Y/4xpeWuckpVC9VJE1ZZPQtA="),
					func(*terraform.State) error {
						assert.False(isFileExists("my-app.zip"))
						assert.False(isFileExists("app/exec.txt"))
						return nil
					},
				),
			},
		},
	})
}

func TestContents_basic(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	os.Mkdir("app", 0755)

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			// Step 1 =====================================================
			{
				Config: `
					resource "lambdazip_file" "my_app" {
						base_dir = "app"
						contents = {
							"app/hello.rb"     = "puts 'world'"
							"app/world.rb"     = "puts 'hello'"
							"app/README.md"    = "# hello.rb"
							"app/lib/const.rb" = "A = 100"
						}
						output = "my-app.zip"

						triggers = {
							hello_rb = "06db2c7a260efaf6e2e3f4c635c83506f1f40f6d3898e0e6025e3e55f44ddebe"
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base_dir", "app"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.#", "0"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "contents.%", "4"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", `contents.app/hello.rb`, "puts 'world'"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "contents.app/hello.rb", "puts 'world'"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "contents.app/README.md", "# hello.rb"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "contents.app/lib/const.rb", "A = 100"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "output", "my-app.zip"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "triggers.hello_rb", "06db2c7a260efaf6e2e3f4c635c83506f1f40f6d3898e0e6025e3e55f44ddebe"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base64sha256", "3RvIR+rYfQjlytlx/hGPe2drWNSu59499c3uAKz6Kh4="),
					func(*terraform.State) error {
						buf, err := os.ReadFile("my-app.zip")
						require.NoError(err)
						assert.Equal("3RvIR+rYfQjlytlx/hGPe2drWNSu59499c3uAKz6Kh4=", base64Sha256(buf))
						list, err := listZip(buf)
						require.NoError(err)
						assert.Equal([]string{"app/README.md", "app/hello.rb", "app/lib/const.rb", "app/world.rb"}, list)
						return nil
					},
				),
			},
			// Step 2 =====================================================
			{
				Config: `
					resource "lambdazip_file" "my_app" {
						base_dir = "app"
						contents = {
							"app/hello.rb"     = "puts 'world'"
							"app/world.rb"     = "puts 'hello'"
							"app/README.md"    = "# hello.rb"
							"app/lib/const.rb" = "A = 100"
						}
						output = "my-app.zip"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base_dir", "app"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.#", "0"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "contents.%", "4"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", `contents.app/hello.rb`, "puts 'world'"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "contents.app/hello.rb", "puts 'world'"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "contents.app/README.md", "# hello.rb"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "contents.app/lib/const.rb", "A = 100"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "output", "my-app.zip"),
					resource.TestCheckNoResourceAttr("lambdazip_file.my_app", "triggers"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base64sha256", "3RvIR+rYfQjlytlx/hGPe2drWNSu59499c3uAKz6Kh4="),
					func(*terraform.State) error {
						buf, err := os.ReadFile("my-app.zip")
						require.NoError(err)
						assert.Equal("3RvIR+rYfQjlytlx/hGPe2drWNSu59499c3uAKz6Kh4=", base64Sha256(buf))
						list, err := listZip(buf)
						require.NoError(err)
						assert.Equal([]string{"app/README.md", "app/hello.rb", "app/lib/const.rb", "app/world.rb"}, list)
						return nil
					},
				),
			},
		},
	})
}

func TestFiles_tempDir(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cwd, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(cwd)

	os.Mkdir("app", 0755)
	os.Mkdir("app/lib", 0755)
	os.WriteFile("app/hello.rb", []byte("puts 'world'"), 0755)
	os.WriteFile("app/world.rb", []byte("puts 'hello'"), 0755)
	os.WriteFile("app/README.md", []byte("# hello.rb"), 0644)
	os.WriteFile("app/lib/const.rb", []byte("A = 100"), 0644)

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			// Step 1 =====================================================
			{
				Config: `
					resource "lambdazip_file" "my_app" {
						base_dir      = "app"
						sources       = ["**/*.rb", "**/*.txt"]
						excludes      = [".*", "README.md"]
						output        = "my-app.zip"
						before_create = "touch exec.txt"
						use_temp_dir  = true

						triggers = {
							hello_rb = filesha256("app/hello.rb"),
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base_dir", "app"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.#", "2"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.0", "**/*.rb"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.1", "**/*.txt"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.#", "2"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.0", ".*"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.1", "README.md"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "output", "my-app.zip"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "before_create", "touch exec.txt"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "use_temp_dir", "true"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "triggers.hello_rb", "06db2c7a260efaf6e2e3f4c635c83506f1f40f6d3898e0e6025e3e55f44ddebe"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base64sha256", "3mKBu6qleIlPZxj0rzLcS52hB1p0K2L76H9MZb+h2kk="),
					func(*terraform.State) error {
						buf, err := os.ReadFile("my-app.zip")
						require.NoError(err)
						assert.Equal("3mKBu6qleIlPZxj0rzLcS52hB1p0K2L76H9MZb+h2kk=", base64Sha256(buf))
						assert.False(isFileExists("app/exec.txt"))
						list, err := listZip(buf)
						require.NoError(err)
						assert.Equal([]string{"exec.txt", "hello.rb", "lib/const.rb", "world.rb"}, list)
						return nil
					},
				),
			},
			// Step 2 =====================================================
			{
				Config: `
					resource "lambdazip_file" "my_app" {
						base_dir      = "app"
						sources       = ["**/*.rb", "**/*.txt"]
						excludes      = [".*", "README.md"]
						output        = "my-app.zip"
						before_create = "touch exec.txt"
						use_temp_dir  = true

						triggers = {
							hello_rb = filesha256("app/hello.rb"),
						}
					}
				`,
				PreConfig: func() {
					err := os.Remove("my-app.zip")
					require.NoError(err)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base_dir", "app"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.#", "2"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.0", "**/*.rb"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.1", "**/*.txt"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.#", "2"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.0", ".*"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.1", "README.md"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "output", "my-app.zip"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "before_create", "touch exec.txt"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "use_temp_dir", "true"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "triggers.hello_rb", "06db2c7a260efaf6e2e3f4c635c83506f1f40f6d3898e0e6025e3e55f44ddebe"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base64sha256", "3mKBu6qleIlPZxj0rzLcS52hB1p0K2L76H9MZb+h2kk="),
					func(*terraform.State) error {
						assert.False(isFileExists("my-app.zip"))
						assert.False(isFileExists("app/exec.txt"))
						return nil
					},
				),
			},
			// Step 3 =====================================================
			{
				Config: `
					resource "lambdazip_file" "my_app" {
						base_dir      = "app"
						sources       = ["**/*.rb", "**/*.txt"]
						excludes      = [".*", "README.md"]
						output        = "my-app.zip"
						before_create = "touch exec.txt"
						use_temp_dir  = true

						triggers = {
							hello_rb = filesha256("app/hello.rb"),
						}
					}
				`,
				PreConfig: func() {
					err := os.WriteFile("app/hello.rb", []byte("print 'world'"), 0755)
					require.NoError(err)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base_dir", "app"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.#", "2"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.0", "**/*.rb"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.1", "**/*.txt"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.#", "2"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.0", ".*"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.1", "README.md"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "output", "my-app.zip"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "before_create", "touch exec.txt"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "use_temp_dir", "true"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "triggers.hello_rb", "6740287d0049734d6fe501a11d8572ba1befdc690d08d891db539d2f8a9d7273"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base64sha256", "mR/NIFGb1EjnvdQWtPZTA3jnrDUGo/rPX++UGs1G9XI="),
					func(*terraform.State) error {
						buf, err := os.ReadFile("my-app.zip")
						require.NoError(err)
						assert.Equal("mR/NIFGb1EjnvdQWtPZTA3jnrDUGo/rPX++UGs1G9XI=", base64Sha256(buf))
						assert.False(isFileExists("app/exec.txt"))
						list, err := listZip(buf)
						require.NoError(err)
						assert.Equal([]string{"exec.txt", "hello.rb", "lib/const.rb", "world.rb"}, list)
						return nil
					},
				),
			},
			// Step 4 =====================================================
			{
				Config: `
					resource "lambdazip_file" "my_app" {
						base_dir      = "app"
						sources       = ["**/*.rb", "**/*.txt"]
						excludes      = [".*", "README.md"]
						output        = "my-app.zip"
						before_create = "touch exec.txt"
						use_temp_dir  = true

						triggers = {
							hello_rb = filesha256("app/hello.rb"),
						}
					}
				`,
				PreConfig: func() {
					err := os.Remove("my-app.zip")
					require.NoError(err)
					err = os.WriteFile("app/world.rb", []byte("print 'hello'"), 0755)
					require.NoError(err)
					err = os.Remove("app/lib/const.rb")
					require.NoError(err)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base_dir", "app"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.#", "2"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.0", "**/*.rb"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "sources.1", "**/*.txt"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.#", "2"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.0", ".*"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "excludes.1", "README.md"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "output", "my-app.zip"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "before_create", "touch exec.txt"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "use_temp_dir", "true"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "triggers.hello_rb", "6740287d0049734d6fe501a11d8572ba1befdc690d08d891db539d2f8a9d7273"),
					resource.TestCheckResourceAttr("lambdazip_file.my_app", "base64sha256", "mR/NIFGb1EjnvdQWtPZTA3jnrDUGo/rPX++UGs1G9XI="),
					func(*terraform.State) error {
						assert.False(isFileExists("my-app.zip"))
						assert.False(isFileExists("app/exec.txt"))
						return nil
					},
				),
			},
		},
	})
}
