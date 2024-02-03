package lambdazip_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestFilesSha256_basic(t *testing.T) {
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
					data "lambdazip_files_sha256" "trigger" {
						files = ["**"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.lambdazip_files_sha256.trigger", "map.%", "4"),
					resource.TestCheckResourceAttr("data.lambdazip_files_sha256.trigger", "map.app/hello.rb", "06db2c7a260efaf6e2e3f4c635c83506f1f40f6d3898e0e6025e3e55f44ddebe"),
					resource.TestCheckResourceAttr("data.lambdazip_files_sha256.trigger", "map.app/world.rb", "293c10e07909b3a823d7d2ba87c6cdf7400c9ed70132c2c952d7c8147d945a74"),
					resource.TestCheckResourceAttr("data.lambdazip_files_sha256.trigger", "map.app/README.md", "29200c6da7d08c5115ad63fe7b9c542e5d8e8cf8a185f5cd49d2ce71fcde8d75"),
					resource.TestCheckResourceAttr("data.lambdazip_files_sha256.trigger", "map.app/lib/const.rb", "62ed3c7896eb965afcfabafe23828a526fcc4fdc8c9e43ed65f3ffecf140036f"),
				),
			},
			// Step 2 =====================================================
			{
				Config: `
					data "lambdazip_files_sha256" "trigger" {
						files = ["**/*.rb"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.lambdazip_files_sha256.trigger", "map.%", "3"),
					resource.TestCheckResourceAttr("data.lambdazip_files_sha256.trigger", "map.app/hello.rb", "06db2c7a260efaf6e2e3f4c635c83506f1f40f6d3898e0e6025e3e55f44ddebe"),
					resource.TestCheckResourceAttr("data.lambdazip_files_sha256.trigger", "map.app/world.rb", "293c10e07909b3a823d7d2ba87c6cdf7400c9ed70132c2c952d7c8147d945a74"),
					resource.TestCheckResourceAttr("data.lambdazip_files_sha256.trigger", "map.app/lib/const.rb", "62ed3c7896eb965afcfabafe23828a526fcc4fdc8c9e43ed65f3ffecf140036f"),
				),
			},
			// Step 3 =====================================================
			{
				Config: `
					data "lambdazip_files_sha256" "trigger" {
						files = ["app/lib/*.rb", "app/*.md"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.lambdazip_files_sha256.trigger", "map.%", "2"),
					resource.TestCheckResourceAttr("data.lambdazip_files_sha256.trigger", "map.app/README.md", "29200c6da7d08c5115ad63fe7b9c542e5d8e8cf8a185f5cd49d2ce71fcde8d75"),
					resource.TestCheckResourceAttr("data.lambdazip_files_sha256.trigger", "map.app/lib/const.rb", "62ed3c7896eb965afcfabafe23828a526fcc4fdc8c9e43ed65f3ffecf140036f"),
				),
			},
			// Step 5 =====================================================
			{
				Config: `
					data "lambdazip_files_sha256" "trigger" {
						files = ["**"]
						excludes = ["app/*.md", "*/*.rb"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.lambdazip_files_sha256.trigger", "map.%", "1"),
					resource.TestCheckResourceAttr("data.lambdazip_files_sha256.trigger", "map.app/lib/const.rb", "62ed3c7896eb965afcfabafe23828a526fcc4fdc8c9e43ed65f3ffecf140036f"),
				),
			},
			// Step 6 =====================================================
			{
				Config: `
					data "lambdazip_files_sha256" "trigger" {
						files = ["**/*rb"]
						excludes = ["app/world.*"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.lambdazip_files_sha256.trigger", "map.%", "2"),
					resource.TestCheckResourceAttr("data.lambdazip_files_sha256.trigger", "map.app/hello.rb", "06db2c7a260efaf6e2e3f4c635c83506f1f40f6d3898e0e6025e3e55f44ddebe"),
					resource.TestCheckResourceAttr("data.lambdazip_files_sha256.trigger", "map.app/lib/const.rb", "62ed3c7896eb965afcfabafe23828a526fcc4fdc8c9e43ed65f3ffecf140036f"),
				),
			},
		},
	})
}
