package lambdazip

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
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
			"sources": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
				AtLeastOneOf: []string{
					"sources",
					"contents",
				},
			},
			"contents": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
				AtLeastOneOf: []string{
					"sources",
					"contents",
				},
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
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
			"base64sha256": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"use_temp_dir": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"compression_level": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
				ForceNew: true,
			},
		},
	}
}

func createFile(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	output := d.Get("output").(string)
	baseDir := d.Get("base_dir").(string)
	useTempDir := d.Get("use_temp_dir").(bool)
	compressionLevel := d.Get("compression_level").(int)
	cwd, err := os.Getwd()

	if err != nil {
		return diag.FromErr(err)
	}

	if !strings.HasPrefix(output, "/") {
		output = filepath.Join(cwd, output)
	}

	if baseDir != "" {
		err = os.Chdir(baseDir)

		if err != nil {
			return diag.FromErr(err)
		}

		defer os.Chdir(cwd) //nolint:errcheck
	}

	if useTempDir {
		tempDir, err := os.MkdirTemp("", "lambdazip")

		if err != nil {
			return diag.FromErr(err)
		}

		defer os.RemoveAll(tempDir)
		err = os.CopyFS(tempDir, os.DirFS("."))

		if err != nil {
			return diag.FromErr(err)
		}

		err = os.Chdir(tempDir)

		if err != nil {
			return diag.FromErr(err)
		}

		defer os.Chdir(cwd) //nolint:errcheck
	}

	sources := []string{}

	if patterns, ok := d.GetOk("sources"); ok {
		for _, pat := range patterns.([]any) {
			sources = append(sources, pat.(string))
		}

		if len(sources) == 0 {
			return diag.Errorf(`The attribute "sources" is required, but the list was empty.`)
		}

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

		sources, err = glob.Glob(sources, excludes)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	contents := map[string]string{}

	if dataMap, ok := d.GetOk("contents"); ok {
		for name, data := range dataMap.(map[string]any) {
			contents[name] = data.(string)
		}
	}

	err = zip.ZipFile(sources, contents, output, compressionLevel)

	if err != nil {
		return diag.FromErr(err)
	}

	base64Sha256, err := hash.Base64Sha256(output)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("base64sha256", base64Sha256) //nolint:errcheck
	d.SetId(id.UniqueId())

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
