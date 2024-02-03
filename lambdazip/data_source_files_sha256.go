package lambdazip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/winebarrel/terraform-provider-lambdazip/internal/glob"
	"github.com/winebarrel/terraform-provider-lambdazip/internal/hash"
)

func dataSourceFilesSha256() *schema.Resource {
	return &schema.Resource{
		ReadContext: readExpr,
		Schema: map[string]*schema.Schema{
			"files": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"excludes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"map": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func readExpr(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	files := []string{}

	if patterns, ok := d.GetOk("files"); ok {
		for _, pat := range patterns.([]any) {
			files = append(files, pat.(string))
		}
	}

	if len(files) == 0 {
		return diag.Errorf(`The attribute "files" is required, but the list was empty.`)
	}

	excludes := []string{}

	if patterns, ok := d.GetOk("excludes"); ok {
		for _, pat := range patterns.([]any) {
			excludes = append(excludes, pat.(string))
		}
	}

	files, err := glob.Glob(files, excludes)

	if err != nil {
		return diag.FromErr(err)
	}

	m, err := hash.Sha256Map(files)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("map", m) //nolint:errcheck
	d.SetId(id.UniqueId())

	return nil
}
