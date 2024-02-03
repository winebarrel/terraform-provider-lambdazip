package lambdazip

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"lambdazip_file": resourceFile(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"lambdazip_files_sha256": dataSourceFilesSha256(),
		},
	}
}
