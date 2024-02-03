package lambdazip

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/winebarrel/terraform-provider-lambdazip/internal/cmd"
	"github.com/winebarrel/terraform-provider-lambdazip/internal/glob"
	"github.com/winebarrel/terraform-provider-lambdazip/internal/hash"
	"github.com/winebarrel/terraform-provider-lambdazip/internal/zip"
)

func resourceFile() *schema.Resource {
	return &schema.Resource{
		CreateContext: createFile,
		ReadContext:   readFile,
		DeleteContext: deleteFile,
		Schema: map[string]*schema.Schema{
			"base_dir": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"excludes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
			"output": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"before_create": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"triggers": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				ForceNew: true,
			},
			"base64sha256": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func createFile(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	baseDir := d.Get("base_dir").(string)
	cwd, err := os.Getwd()

	if err != nil {
		return diag.FromErr(err)
	}

	if baseDir != "" {
		err = os.Chdir(baseDir)

		if err != nil {
			return diag.FromErr(err)
		}

		defer os.Chdir(cwd) //nolint:errcheck
	}

	source := d.Get("source").(string)
	output := d.Get("output").(string)
	excludes := []string{}

	if patterns, ok := d.GetOk("excludes"); ok {
		for _, pat := range patterns.([]any) {
			excludes = append(excludes, pat.(string))
		}
	}

	if beforeCreate, ok := d.GetOk("before_create"); ok {
		beforeCreateCmd := beforeCreate.(string)
		cmdout, err := cmd.Run(beforeCreateCmd)

		if err != nil {
			errmsg := `"%s" failed - %w`

			if cmdout != "" {
				errmsg += ":\n" + cmdout
			}

			return diag.Errorf(errmsg, beforeCreateCmd, err)
		}
	}

	sources, err := glob.Glob(source, excludes)

	if err != nil {
		return diag.FromErr(err)
	}

	if baseDir != "" && !strings.HasPrefix("/", baseDir) {
		output = filepath.Join(cwd, output)
	}

	err = zip.ZipFile(sources, output)

	if err != nil {
		return diag.FromErr(err)
	}

	base64Sha256, err := hash.Base64Sha256(output)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("base64sha256", base64Sha256) //nolint:errcheck
	d.SetId(output)

	return nil
}

func readFile(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func deleteFile(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	output := d.Get("output").(string)
	os.Remove(output)
	d.SetId("")

	return nil
}
