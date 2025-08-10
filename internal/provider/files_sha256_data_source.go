package provider

import (
	"context"
	"maps"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/winebarrel/terraform-provider-lambdazip/internal/glob"
	"github.com/winebarrel/terraform-provider-lambdazip/internal/hash"
)

var _ datasource.DataSourceWithConfigValidators = &FilesSha256DataSource{}

func NewFilesSha256DataSource() datasource.DataSource {
	return &FilesSha256DataSource{}
}

type FilesSha256DataSource struct {
}

type FilesSha256DataSourceModel struct {
	Files         []types.String `tfsdk:"files"`
	Contents      types.Map      `tfsdk:"contents"`
	Excludes      []types.String `tfsdk:"excludes"`
	Map           types.Map      `tfsdk:"map"`
	AllowNotExist types.Bool     `tfsdk:"allow_not_exist"`
}

func (d *FilesSha256DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_files_sha256"
}

func (d *FilesSha256DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"files": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.NoNullValues(),
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(0)),
				},
			},
			"contents": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.Map{
					mapvalidator.NoNullValues(),
					mapvalidator.SizeAtLeast(1),
					mapvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(0)),
				},
			},
			"excludes": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.NoNullValues(),
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(0)),
				},
			},
			"map": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"allow_not_exist": schema.BoolAttribute{
				Optional: true,
			},
		},
	}
}

func (d *FilesSha256DataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.AtLeastOneOf(
			path.MatchRoot("files"),
			path.MatchRoot("contents"),
		),
	}
}

func (d *FilesSha256DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FilesSha256DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	mFiles := map[string]string{}

	if len(data.Files) >= 1 {
		files := []string{}

		for _, f := range data.Files {
			files = append(files, f.ValueString())
		}

		excludes := []string{}

		for _, e := range data.Excludes {
			excludes = append(excludes, e.ValueString())
		}

		globOpts := []doublestar.GlobOption{}

		if !data.AllowNotExist.ValueBool() {
			globOpts = append(globOpts, doublestar.WithFailOnPatternNotExist())
		}

		files, err := glob.Glob(files, excludes, globOpts...)

		if err != nil {
			resp.Diagnostics.AddError("Failed to glob files", err.Error())
			return
		}

		mFiles, err = hash.Sha256Map(files)

		if err != nil {
			resp.Diagnostics.AddError("Failed to calculate sha256sum", err.Error())
			return
		}
	}

	mContents := map[string]string{}

	if len(data.Contents.Elements()) >= 1 {
		elements := make(map[string]types.String, len(data.Contents.Elements()))
		data.Contents.ElementsAs(ctx, &elements, false)
		dataByName := map[string]string{}

		for name, data := range elements {
			dataByName[name] = data.ValueString()
		}

		mContents = hash.ContentsSha256Map(dataByName)
	}

	maps.Copy(mFiles, mContents)
	m, diags := types.MapValueFrom(ctx, types.StringType, mFiles)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Map = m
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
