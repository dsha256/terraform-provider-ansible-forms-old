package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-ansibleforms/internal/interfaces"
	"terraform-provider-ansibleforms/internal/utils"
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
	CxProfileName types.String  `tfsdk:"cx_profile_name"`
	ID            types.Int64   `tfsdk:"id"`
	Status        types.String  `tfsdk:"status"`
	LastUpdated   types.String  `tfsdk:"last_updated"`
	FormName      types.String  `tfsdk:"form_name"`
	ExtraVars     types.MapType `tfsdk:"extravars"`
}

// JobResourceModelCredentials ...
type JobResourceModelCredentials struct {
	OntapCred types.String `tfsdk:"ontap_cred"`
	BindCred  types.String `tfsdk:"bind_cred"`
}

// Metadata returns the resource type name.
func (r *JobResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	//resp.TypeName = req.ProviderTypeName + "_" + r.config.name
	resp.TypeName = "ansibleforms_job"
}

// Schema defines the schema for the resource.
func (r *JobResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Job resource",

		Attributes: map[string]schema.Attribute{
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
			"form_name": schema.StringAttribute{
				Description: "Form Name.",
				Optional:    true,
			},
			"extravars": schema.MapAttribute{
				Description: "Extra Vars.",
				Required:    true,
				ElementType: types.StringType,
			},
			"credentials": schema.SingleNestedAttribute{
				Description: "Credentials",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"ontap_cred": schema.StringAttribute{
						Description: "OnTap Credentials.",
						Required:    true,
					},
					"bind_cred": schema.StringAttribute{
						Description: "Bind Credentials.",
						Required:    true,
					},
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
	//// Retrieve values from plan
	//var plan JobResourceModel
	//diags := req.Plan.Get(ctx, &plan)
	//resp.Diagnostics.Append(diags...)
	//if resp.Diagnostics.HasError() {
	//	return
	//}
	//
	//// Generate API request body from plan
	//var createjob AnsibleFormsCreateJob
	//createjob.FormName = string(plan.FormName.ValueString())
	//
	//createjob.ExtraVars.AccountID = plan.ExtraVars.AccountID.ValueString()
	//createjob.ExtraVars.Dataclass = plan.ExtraVars.Dataclass.ValueString()
	//createjob.ExtraVars.Env = plan.ExtraVars.Env.ValueString()
	//createjob.ExtraVars.Exposure = plan.ExtraVars.Exposure.ValueString()
	//createjob.ExtraVars.Opco = plan.ExtraVars.Opco.ValueString()
	//createjob.ExtraVars.ProtectionRequired = plan.ExtraVars.ProtectionRequired.ValueString()
	//createjob.ExtraVars.Region = plan.ExtraVars.Region.ValueString()
	//createjob.ExtraVars.ShareName = plan.ExtraVars.ShareName.ValueString()
	//createjob.ExtraVars.Size = plan.ExtraVars.Size.ValueString()
	//createjob.ExtraVars.State = "present"
	//createjob.ExtraVars.SvmName = plan.ExtraVars.SvmName.ValueString()
	//
	//createjob.Credentials.BindCred = plan.Credentials.BindCred.ValueString()
	//createjob.Credentials.OntapCred = plan.Credentials.OntapCred.ValueString()
	//
	//// Create new job
	//job, err := r.client.CreateJob(ctx, createjob)
	//if err != nil {
	//	resp.Diagnostics.AddError(
	//		"Error creating job",
	//		"Could not create job, unexpected error: "+err.Error(),
	//	)
	//	return
	//}
	//
	//// Map response body to schema and populate Computed attribute values
	//plan.ID = types.StringValue(strconv.Itoa(job.Data.Output.ID))
	//plan.Status = types.StringValue(job.Status)
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	//
	//// Set state to fully populated data
	//diags = resp.State.Set(ctx, plan)
	//resp.Diagnostics.Append(diags...)
	//if resp.Diagnostics.HasError() {
	//	return
	//}
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
		// NOTE: error reporting done inside NewClient
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("read a job resource: %#v", data))

	var job *interfaces.JobGetDataSourceModel
	job, err = interfaces.GetJobById(errorHandler, *client, data.ID.String())
	if err != nil {
		return
	}
	if job == nil {
		errorHandler.MakeAndReportError("No Svm found", "No SVM found")
	}

	//data.FormName = types.StringValue(job.Form)
	//data.ID = types.Int64Value(job.Id.ValueInt64())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	//
	//
	//
	//
	//

	//// Get current state
	//var state JobResourceModel
	//diags := req.State.Get(ctx, &state)
	//resp.Diagnostics.Append(diags...)
	//if resp.Diagnostics.HasError() {
	//	return
	//}
	//
	//// Get refreshed job value from AnsibleForms
	//job, err := r.client.GetJob(ctx, state.ID.ValueString())
	//if err != nil {
	//	resp.Diagnostics.AddError(
	//		"Error Reading AnsibleForms Job",
	//		"Could not read AnsibleForms job ID "+state.ID.ValueString()+": "+err.Error(),
	//	)
	//	return
	//}
	//
	//// Overwrite items with refreshed state
	//state.FormName = types.StringValue(job.Data.Form)
	//state.Status = types.StringValue(job.Status)
	//
	//// Set refreshed state
	//diags = resp.State.Set(ctx, &state)
	//resp.Diagnostics.Append(diags...)
	//if resp.Diagnostics.HasError() {
	//	return
	//}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *JobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *JobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//// Retrieve values from state
	//var state JobResourceModel
	//diags := req.State.Get(ctx, &state)
	//resp.Diagnostics.Append(diags...)
	//if resp.Diagnostics.HasError() {
	//	return
	//}
	//
	//// Generate API request body from plan
	//var createjob AnsibleFormsCreateJob
	//createjob.FormName = state.FormName.ValueString()
	//
	//createjob.ExtraVars.AccountID = string(state.ExtraVars.AccountID.ValueString())
	//createjob.ExtraVars.Dataclass = string(state.ExtraVars.Dataclass.ValueString())
	//createjob.ExtraVars.Env = string(state.ExtraVars.Env.ValueString())
	//createjob.ExtraVars.Exposure = string(state.ExtraVars.Exposure.ValueString())
	//createjob.ExtraVars.Opco = string(state.ExtraVars.Opco.ValueString())
	//createjob.ExtraVars.ProtectionRequired = string(state.ExtraVars.ProtectionRequired.ValueString())
	//createjob.ExtraVars.Region = string(state.ExtraVars.Region.ValueString())
	//createjob.ExtraVars.ShareName = string(state.ExtraVars.ShareName.ValueString())
	//createjob.ExtraVars.Size = string(state.ExtraVars.Size.ValueString())
	//createjob.ExtraVars.State = "absent"
	//createjob.ExtraVars.SvmName = string(state.ExtraVars.SvmName.ValueString())
	//
	//createjob.Credentials.BindCred = string(state.Credentials.BindCred.ValueString())
	//createjob.Credentials.OntapCred = string(state.Credentials.OntapCred.ValueString())
	//
	//// Delete existing order
	//// Create new job
	//job, err := r.client.CreateJob(ctx, createjob)
	//if err != nil {
	//	resp.Diagnostics.AddError(
	//		"Error Create Job (State Absent)",
	//		"Could not create job (State Absent), unexpected error: "+err.Error(),
	//	)
	//	return
	//}
	//if string(job.Status) != "success" {
	//	resp.Diagnostics.AddError(
	//		"Error Status Job (State Absent)",
	//		"Create job (State Absent) has status , unexpected error: "+string(job.Status),
	//	)
	//	return
	//}

}
