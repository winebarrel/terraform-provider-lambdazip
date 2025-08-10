//go:build generate

package tools

//go:generate terraform fmt -recursive ../examples/
//go:generate go tool github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-dir .. -provider-name lambdazip
