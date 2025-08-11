package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	cp "github.com/otiai10/copy"
	"github.com/winebarrel/terraform-provider-lambdazip/internal/cmd"
	"github.com/winebarrel/terraform-provider-lambdazip/internal/glob"
	"github.com/winebarrel/terraform-provider-lambdazip/internal/hash"
	"github.com/winebarrel/terraform-provider-lambdazip/internal/zip"
)

var _ resource.ResourceWithConfigValidators = &FileResource{}

func NewFileResource() resource.Resource {
	return &FileResource{}
}

type FileResource struct {
}

type FileResourceModel struct {
	BaseDir          types.String   `tfsdk:"base_dir"`
	Sources          []types.String `tfsdk:"sources"`
	Contents         types.Map      `tfsdk:"contents"`
	Excludes         []types.String `tfsdk:"excludes"`
	Output           types.String   `tfsdk:"output"`
	BeforeCreate     types.String   `tfsdk:"before_create"`
	Triggers         types.Map      `tfsdk:"triggers"`
	Base64sha256     types.String   `tfsdk:"base64sha256"`
	UseTempDir       types.Bool     `tfsdk:"use_temp_dir"`
	CompressionLevel types.Int32    `tfsdk:"compression_level"`
	StripComponents  types.Int32    `tfsdk:"strip_components"`
}

func (r *FileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (r *FileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_dir": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sources": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.NoNullValues(),
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(0)),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
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
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"excludes": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.NoNullValues(),
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"output": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"before_create": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"triggers": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.Map{
					mapvalidator.NoNullValues(),
					mapvalidator.SizeAtLeast(1),
					mapvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				},
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"base64sha256": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"use_temp_dir": schema.BoolAttribute{
				Optional: true,
			},
			"compression_level": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Default:  int32default.StaticInt32(-1),
				Validators: []validator.Int32{
					int32validator.Between(-1, 9),
				},
			},
			"strip_components": schema.Int32Attribute{
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
		},
	}
}

func (d *FileResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("sources"),
			path.MatchRoot("contents"),
		),
	}
}

func (r *FileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	output := plan.Output.ValueString()
	baseDir := plan.BaseDir.ValueString()
	useTempDir := plan.UseTempDir.ValueBool()
	compressionLevel := int(plan.CompressionLevel.ValueInt32())
	stripComponents := int(plan.StripComponents.ValueInt32())
	cwd, err := os.Getwd()

	if err != nil {
		resp.Diagnostics.AddError("Failed to get current directory", err.Error())
		return
	}

	if !strings.HasPrefix(output, "/") {
		output = filepath.Join(cwd, output)
	}

	if baseDir != "" {
		err = os.Chdir(baseDir)

		if err != nil {
			resp.Diagnostics.AddError("Failed to change current working directory", err.Error())
			return
		}

		defer os.Chdir(cwd) //nolint:errcheck
	}

	if useTempDir {
		tempDir, err := os.MkdirTemp("", "lambdazip")

		if err != nil {
			resp.Diagnostics.AddError("Failed to create temporary directory", err.Error())
			return
		}

		defer os.RemoveAll(tempDir)
		err = cp.Copy(".", tempDir)

		if err != nil {
			resp.Diagnostics.AddError("Failed to copy files to temporary directory", err.Error())
			return
		}

		err = os.Chdir(tempDir)

		if err != nil {
			resp.Diagnostics.AddError("Failed to change current working directory", err.Error())
			return
		}

		defer os.Chdir(cwd) //nolint:errcheck
	}

	sources := []string{}

	if len(plan.Sources) >= 1 {
		for _, pat := range plan.Sources {
			sources = append(sources, pat.ValueString())
		}

		excludes := []string{}

		for _, pat := range plan.Excludes {
			excludes = append(excludes, pat.ValueString())
		}

		if beforeCreate := plan.BeforeCreate.ValueString(); beforeCreate != "" {
			cmdout, err := cmd.Run(beforeCreate)

			if err != nil {
				cmdout = strings.TrimSpace(cmdout)

				if cmdout == "" {
					cmdout = "(empty)"
				}

				summary := fmt.Sprintf("Failed to run `%s`", beforeCreate)
				detail := fmt.Sprintf("%s\noutput: %s", err, cmdout)
				resp.Diagnostics.AddError(summary, detail)
				return
			}
		}

		sources, err = glob.Glob(sources, excludes)

		if err != nil {
			resp.Diagnostics.AddError("Failed to glob files", err.Error())
			return
		}
	}

	contents := map[string]string{}

	if len(plan.Contents.Elements()) >= 1 {
		elements := make(map[string]types.String, len(plan.Contents.Elements()))
		plan.Contents.ElementsAs(ctx, &elements, false)

		for name, data := range elements {
			contents[name] = data.ValueString()
		}
	}

	err = zip.ZipFile(sources, contents, output, compressionLevel, stripComponents)

	if err != nil {
		resp.Diagnostics.AddError("Failed to zip files", err.Error())
		return
	}

	base64sha256, err := hash.Base64Sha256(output)

	if err != nil {
		resp.Diagnostics.AddError("Failed to calculate sha256sum", err.Error())
		return
	}

	plan.Base64sha256 = types.StringValue(base64sha256)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *FileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *FileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}
