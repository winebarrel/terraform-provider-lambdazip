package lambdazip

import (
	"context"

	"github.com/bmatcuk/doublestar/v4"
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
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				AtLeastOneOf: []string{
					"files",
					"contents",
				},
			},
			"contents": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				AtLeastOneOf: []string{
					"files",
					"contents",
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
			"allow_not_exist": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func readExpr(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	mFiles := map[string]string{}
	mContents := map[string]string{}

	if patterns, ok := d.GetOk("files"); ok {
		files := []string{}

		for _, pat := range patterns.([]any) {
			files = append(files, pat.(string))
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

		globOpts := []doublestar.GlobOption{}

		if !d.Get("allow_not_exist").(bool) {
			globOpts = append(globOpts, doublestar.WithFailOnPatternNotExist())
		}

		files, err := glob.Glob(files, excludes, globOpts...)

		if err != nil {
			return diag.FromErr(err)
		}

		mFiles, err = hash.Sha256Map(files)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	if dataMap, ok := d.GetOk("contents"); ok {
		dataByName := map[string]string{}

		for name, data := range dataMap.(map[string]any) {
			dataByName[name] = data.(string)
		}

		mContents = hash.ContentsSha256Map(dataByName)
	}

	for name, hash := range mContents {
		mFiles[name] = hash
	}

	d.Set("map", mFiles) //nolint:errcheck
	d.SetId(id.UniqueId())

	return nil
}
