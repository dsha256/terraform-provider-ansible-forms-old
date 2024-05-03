package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-ansible-forms/internal/interfaces"
	"terraform-provider-ansible-forms/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &JobResource{}
	_ resource.ResourceWithConfigure = &JobResource{}
)

// NewJobResource is a helper function to simplify the provider implementation.
func NewJobResource() resource.Resource {
	return &JobResource{
		config: resourceOrDataSourceConfig{
			name: "job_resource",
		},
	}
}

// JobResource is the resource implementation.
type JobResource struct {
	config resourceOrDataSourceConfig
}

// JobResourceModel maps the resource schema data.
type JobResourceModel struct {
	CxProfileName types.String `tfsdk:"cx_profile_name"`
	ID            types.String `tfsdk:"id"`
	LastUpdated   types.String `tfsdk:"last_updated"`
	FormName      types.String `tfsdk:"form_name"`
	Status        types.String `tfsdk:"status"`
	Extravars     types.Map    `tfsdk:"extravars"`
	Credentials   types.Map    `tfsdk:"credentials"`
}

// JobResourceModelCredentials ...
type JobResourceModelCredentials struct {
	OntapCred types.String `tfsdk:"ontap_cred"`
	BindCred  types.String `tfsdk:"bind_cred"`
}

// Metadata returns the resource type name.
func (r *JobResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.config.name
}

// Schema defines the schema for the resource.
func (r *JobResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Job resource",

		Attributes: map[string]schema.Attribute{
			"cx_profile_name": schema.StringAttribute{
				MarkdownDescription: "Connection profile name",
				Required:            true,
			},
			"form_name": schema.StringAttribute{
				Description: "Form Name.",
				Required:    true,
			},
			"extravars": schema.MapAttribute{
				Description: "Extra Vars.",
				Required:    true,
				ElementType: types.StringType,
			},
			"credentials": schema.MapAttribute{
				Description: "Extra Vars.",
				Required:    true,
				ElementType: types.StringType,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *JobResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	config, ok := req.ProviderData.(Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected  Resource Configure Type",
			fmt.Sprintf("Expected Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
	r.config.providerConfig = config
}

// Create a new resource.
func (r *JobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *JobResourceModel
	errorHandler := utils.NewErrorHandler(ctx, &resp.Diagnostics)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, "error getting req plan")
		return
	}

	var request interfaces.JobResourceModel
	request.Form = data.FormName.ValueString()
	//request.Extravars = data.Extravars.Elements()

	client, err := getRestClient(errorHandler, r.config, data.CxProfileName)
	if err != nil {
		// error reporting done inside NewClient
		return
	}

	job, err := interfaces.CreateJob(errorHandler, *client, request)
	if err != nil {
		tflog.Debug(ctx, "err creating a resource", map[string]interface{}{"err": err})
		return
	}

	data.ID = types.StringValue(strconv.FormatInt(job.Data.ID, 10))
	data.Status = types.StringValue(job.Data.Status)
	data.LastUpdated = types.StringValue(time.Now().UTC().Format(time.RFC3339))

	tflog.Debug(ctx, "JOB ID", map[string]interface{}{"ID": job.Data.ID, "DATA": data})

	tflog.Trace(ctx, "created a resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read resource information.
func (r *JobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *JobResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	errorHandler := utils.NewErrorHandler(ctx, &resp.Diagnostics)

	client, err := getRestClient(errorHandler, r.config, data.CxProfileName)
	if err != nil {
		// error reporting done inside NewClient
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("read a job resource: %#v", data))

	var job *interfaces.JobGetDataSourceModel
	if data.ID.ValueString() != "" {
		job, err = interfaces.GetJobById(errorHandler, *client, data.ID.ValueString())
	} else {
		return
	}
	if err != nil {
		return
	}

	if job == nil {
		errorHandler.MakeAndReportError("No Job found", "No JOB found")
	}

	restInfo, err := interfaces.GetJobById(errorHandler, *client, data.ID.String())
	if err != nil {
		// error reporting done inside GetSVMPeer
		return
	}

	data.ID = types.StringValue(strconv.FormatInt(restInfo.Data.ID, 10))
	data.FormName = types.StringValue(restInfo.Data.Form)
	data.Status = types.StringValue(restInfo.Data.Status)
	//data.Extravars = jsonStringToMapValue(ctx, &resp.Diagnostics, restInfo.Data.Extravars)
	//data.Credentials = jsonStringToMapValue(ctx, &resp.Diagnostics, restInfo.Data.Credentials)

	//data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Debug(ctx, fmt.Sprintf("read a data source: %#v", data))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *JobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *JobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
