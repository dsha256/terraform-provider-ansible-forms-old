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
)

// Ensure the implementation satisfies the expected interfaces.
// Ensure the implementation satisfies the expected interfaces.
var (
    _ resource.Resource              = &jobResource{}
    _ resource.ResourceWithConfigure = &jobResource{}
)


// JobResource is the resource implementation.
type jobResource struct {
    client *AnsibleFormsClient
}

// jobResourceModel maps the resource schema data.
type jobResourceModel struct {
    ID          types.String     `tfsdk:"id"`
    Status      types.String 	 `tfsdk:"status"`
    LastUpdated types.String     `tfsdk:"last_updated"`
	FormName 	types.String     `tfsdk:"form_name"`
	ExtraVars   jobResourceModelExtraVars	 	`tfsdk:"extravars"`
	Credentials jobResourceModelCredentials    `tfsdk:"credentials"`
}

type jobResourceModelExtraVars struct {
	Region	types.String 	`tfsdk:"region"`
	Opco	types.String 	`tfsdk:"opco"`
    SvmName types.String 	`tfsdk:"svm_name"`
//	State types.String 		`tfsdk:"state"`
	Exposure types.String 	`tfsdk:"exposure"`
	Env types.String 		`tfsdk:"env"`
	Dataclass types.String 	`tfsdk:"dataclass"`
	ShareName types.String 	`tfsdk:"share_name"`
	AccountID types.String 	`tfsdk:"accountid"`
	Size types.String 		`tfsdk:"size"`
	ProtectionRequired types.String `tfsdk:"protection_required"`
}

type jobResourceModelCredentials struct {     
    OntapCred types.String 	`tfsdk:"ontap_cred"`
    BindCred types.String 	`tfsdk:"bind_cred"`
}

// NewJobrResource is a helper function to simplify the provider implementation.
func NewJobResource() resource.Resource {
    return &jobResource{}
}

// Metadata returns the resource type name.
func (r *jobResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_job"
}

// Schema defines the schema for the resource.
func (r *jobResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
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
                Required: true,
            },
            "extravars": schema.SingleNestedAttribute{
				Description: "ExtraVars",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"region": schema.StringAttribute{
						Description: "Region.",
						Required:    true,
					},
					"opco": schema.StringAttribute{
						Description: "OPCO.",
						Required:    true,
					},
					"svm_name": schema.StringAttribute{
						Description: "SVM Name.",
						Required:    true,
					},
//					"state": schema.StringAttribute{
//						Description: "State.",
//						Required:    true,
//					},
					"exposure": schema.StringAttribute{
						Description: "Exposure.",
						Required:    true,
					},
					"env": schema.StringAttribute{
						Description: "Environment.",
						Required:    true,
					},
					"dataclass": schema.StringAttribute{
						Description: "Data Class.",
						Required:    true,
					},
					"share_name": schema.StringAttribute{
						Description: "Share Name.",
						Required:    true,
					},
					"accountid": schema.StringAttribute{
						Description: "Account ID.",
						Required:    true,
					},
					"size": schema.StringAttribute{
						Description: "Size.",
						Required:    true,
					},
					"protection_required": schema.StringAttribute{
						Description: "Protection Required.",
						Required:    true,
					},
				},
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

// Create a new resource.
func (r *jobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Retrieve values from plan
    var plan jobResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

	// Generate API request body from plan
	var createjob AnsibleFormsCreateJob
	createjob.FormName=string(plan.FormName.ValueString())

	createjob.ExtraVars.AccountID = string(plan.ExtraVars.AccountID.ValueString())
	createjob.ExtraVars.Dataclass = string(plan.ExtraVars.Dataclass.ValueString())
	createjob.ExtraVars.Env = string(plan.ExtraVars.Env.ValueString())
	createjob.ExtraVars.Exposure = string(plan.ExtraVars.Exposure.ValueString())
	createjob.ExtraVars.Opco = string(plan.ExtraVars.Opco.ValueString())
	createjob.ExtraVars.ProtectionRequired = string(plan.ExtraVars.ProtectionRequired.ValueString())
	createjob.ExtraVars.Region = string(plan.ExtraVars.Region.ValueString())
	createjob.ExtraVars.ShareName = string(plan.ExtraVars.ShareName.ValueString())
	createjob.ExtraVars.Size = string(plan.ExtraVars.Size.ValueString())
	createjob.ExtraVars.State = "present"
	createjob.ExtraVars.SvmName = string(plan.ExtraVars.SvmName.ValueString())
	
	createjob.Credentials.BindCred=string(plan.Credentials.BindCred.ValueString())
	createjob.Credentials.OntapCred=string(plan.Credentials.OntapCred.ValueString())

    // Create new job
    job, err := r.client.CreateJob(ctx, createjob)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error creating job",
            "Could not create job, unexpected error: "+err.Error(),
        )
        return
    }

    // Map response body to schema and populate Computed attribute values
    plan.ID = types.StringValue(strconv.Itoa(job.Data.Output.ID))
	plan.Status = types.StringValue(job.Status)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))


    // Set state to fully populated data
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Read resource information.
func (r *jobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
		var state jobResourceModel
		diags := req.State.Get(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	
		// Get refreshed job value from AnsibleForms
		job, err := r.client.GetJob(ctx, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading AnsibleForms Job",
				"Could not read AnsibleForms job ID "+state.ID.ValueString()+": "+err.Error(),
			)
			return
		}
	
		// Overwrite items with refreshed state
		state.FormName = types.StringValue(job.Data.Form)
		state.Status = types.StringValue(job.Status)
	
		// Set refreshed state
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	
// Update updates the resource and sets the updated Terraform state on success.
func (r *jobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *jobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
// Retrieve values from state
var state jobResourceModel
diags := req.State.Get(ctx, &state)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
	return
}

// Generate API request body from plan
var createjob AnsibleFormsCreateJob
createjob.FormName=string(state.FormName.ValueString())

createjob.ExtraVars.AccountID = string(state.ExtraVars.AccountID.ValueString())
createjob.ExtraVars.Dataclass = string(state.ExtraVars.Dataclass.ValueString())
createjob.ExtraVars.Env = string(state.ExtraVars.Env.ValueString())
createjob.ExtraVars.Exposure = string(state.ExtraVars.Exposure.ValueString())
createjob.ExtraVars.Opco = string(state.ExtraVars.Opco.ValueString())
createjob.ExtraVars.ProtectionRequired = string(state.ExtraVars.ProtectionRequired.ValueString())
createjob.ExtraVars.Region = string(state.ExtraVars.Region.ValueString())
createjob.ExtraVars.ShareName = string(state.ExtraVars.ShareName.ValueString())
createjob.ExtraVars.Size = string(state.ExtraVars.Size.ValueString())
createjob.ExtraVars.State = "absent"
createjob.ExtraVars.SvmName = string(state.ExtraVars.SvmName.ValueString())

createjob.Credentials.BindCred=string(state.Credentials.BindCred.ValueString())
createjob.Credentials.OntapCred=string(state.Credentials.OntapCred.ValueString())

// Delete existing order
// Create new job
job, err := r.client.CreateJob(ctx, createjob)
if err != nil {
	resp.Diagnostics.AddError(
		"Error Create Job (State Absent)",
		"Could not create job (State Absent), unexpected error: "+err.Error(),
	)
	return
}
if string(job.Status) != "success" {
	resp.Diagnostics.AddError(
		"Error Status Job (State Absent)",
		"Create job (State Absent) has status , unexpected error: " + string(job.Status),
	)
	return
}

}

// Configure adds the provider configured client to the resource.
func (r *jobResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*AnsibleFormsClient)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Data Source Configure Type",
            fmt.Sprintf("Expected *AnsibleFormsClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}
