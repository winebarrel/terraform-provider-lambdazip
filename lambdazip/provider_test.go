package lambdazip_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/terraform-provider-lambdazip/lambdazip"
)

var (
	testProviders map[string]*schema.Provider
	testProvider  *schema.Provider
)

func init() {
	testProvider = lambdazip.Provider()
	testProviders = map[string]*schema.Provider{
		"lambdazip": testProvider,
	}
}

func TestProvider(t *testing.T) {
	assert := assert.New(t)
	provider := lambdazip.Provider()
	err := provider.InternalValidate()
	assert.NoError(err)
}
