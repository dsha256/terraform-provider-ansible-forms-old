// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &ansibleformsProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ansibleformsProvider{
			version: version,
		}
	}
}

// ansibleformsProvider is the provider implementation.
type ansibleformsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ansibleformsProviderModel maps provider schema data to a Go type.
type ansibleformsProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Metadata returns the provider type name.
func (p *ansibleformsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ansibleforms"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *ansibleformsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with AnsibleForms.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for AnsibleForms API. May also be provided via ANSIBLEFORMS_HOST environment variable.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username for AnsibleForms API. May also be provided via ANSIBLEFORMS_USERNAME environment variable.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for AnsibleForms API. May also be provided via ANSIBLEFORMS_PASSWORD environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *ansibleformsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring AnsibleForms client")

	// Retrieve provider data from configuration
	var config ansibleformsProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown AnsibleForms API Host",
			"The provider cannot create the AnsibleForms API client as there is an unknown configuration value for the AnsibleForms API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ANSIBLEFORMS_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown AnsibleForms API Username",
			"The provider cannot create the AnsibleForms API client as there is an unknown configuration value for the AnsibleForms API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ANSIBLEFORMS_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown AnsibleForms API Password",
			"The provider cannot create the AnsibleForms API client as there is an unknown configuration value for the AnsibleForms API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ANSIBLEFORMS_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("ANSIBLEFORMS_HOST")
	username := os.Getenv("ANSIBLEFORMS_USERNAME")
	password := os.Getenv("ANSIBLEFORMS_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing AnsibleForms API Host",
			"The provider cannot create the AnsibleForms API client as there is a missing or empty value for the AnsibleForms API host. "+
				"Set the host value in the configuration or use the ANSIBLEFORMS_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing AnsibleForms API Username",
			"The provider cannot create the AnsibleForms API client as there is a missing or empty value for the AnsibleForms API username. "+
				"Set the username value in the configuration or use the ANSIBLEFORMS_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing AnsibleForms API Password",
			"The provider cannot create the AnsibleForms API client as there is a missing or empty value for the AnsibleForms API password. "+
				"Set the password value in the configuration or use the ANSIBLEFORMS_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "ansibleforms_host", host)
	ctx = tflog.SetField(ctx, "ansibleforms_username", username)
	ctx = tflog.SetField(ctx, "ansibleforms_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "ansibleforms_password")

	tflog.Debug(ctx, "Creating AnsibleForms client")

	// Create a new AnsibleForms client using the configuration values
	client, err := ansibleformsNewClient(ctx, &host, &username, &password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create AnsibleForms API Client",
			"An unexpected error occurred when creating the AnsibleForms API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"AnsibleForms Client Error: "+err.Error(),
		)
		return
	}

	// Make the AnsibleForms client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured AnsibleForms client", map[string]any{"success": true})

}

// DataSources defines the data sources implemented in the provider.
func (p *ansibleformsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *ansibleformsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewJobResource,
	}
}
