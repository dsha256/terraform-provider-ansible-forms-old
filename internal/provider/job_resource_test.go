package provider

import (
  "testing"

  "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJobResource(t *testing.T) {
  resource.Test(t, resource.TestCase{
    ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
    Steps: []resource.TestStep{
      // Create and Read testing
      {
        Config: providerConfig + `
resource "ansibleforms_job" "test" {
	form_name = "AXA Share Create"
    extravars  = {
     region = "myregion"
     opco   = "myopco"
     svm_name= "mysvm_name"
     exposure = "myexposure"
     env    = "myenv"
     dataclass = "mydataclass"
     share_name = "myshare_name"
     accountid = "myaccountid"
     size = "mysize"
     protection_required = "myprotection_required"
  }
  credentials = {
    ontap_cred = "myontap_cred"
    bind_cred = "mybind_cred"
  }
}
`,
        Check: resource.ComposeAggregateTestCheckFunc(
          // Verify Ansible Forms Name
          resource.TestCheckResourceAttr("ansibleforms_job.test", "form_name", "AXA Share Create"),
          // Verify Ansible Job ExtraVars & Credentials Attributes
          //resource.TestCheckResourceAttr("ansibleforms_job.test", "extravars.region", "myregion"),
          //resource.TestCheckResourceAttr("ansibleforms_job.test", "credentials.ontap_cred", "myontap_cred"),
		  // Verify dynamic values have any value set in the state.
          //resource.TestCheckResourceAttrSet("ansibleforms_job.test", "id"),
          //resource.TestCheckResourceAttrSet("ansibleforms_job.test", "last_updated"),
		  //resource.TestCheckResourceAttrSet("ansibleforms_job.test", "status"),
        ),
      },
      // Update and Read testing
      {
        Config: providerConfig + `
resource "ansibleforms_job" "test" {
    form_name = "AXA Share Create"
	extravars  = {
		region = "myregion"
		opco   = "myopco"
		svm_name= "mysvm_name"
		exposure = "myexposure"
		env    = "myenv"
		dataclass = "mydataclass"
		share_name = "myshare_name"
		accountid = "myaccountid"
		size = "mysize"
		protection_required = "myprotection_required"
	}
	credentials = {
		ontap_cred = "myontap_cred"
		bind_cred = "mybind_cred"
	}
}`,
        Check: resource.ComposeAggregateTestCheckFunc(
           // Verify Ansible Forms Name
           resource.TestCheckResourceAttr("ansibleforms_job.test", "form_name", "AXA Share Create"),
           // Verify Ansible Job ExtraVars & Credentials Attributes
		   resource.TestCheckResourceAttr("ansibleforms_job.test", "extravars.region", "myregion"),
		   resource.TestCheckResourceAttr("ansibleforms_job.test", "credentials.ontap_cred", "myontap_cred"),
		   // Verify dynamic values have any value set in the state.
		   resource.TestCheckResourceAttrSet("ansibleforms_job.test", "id"),
		   resource.TestCheckResourceAttrSet("ansibleforms_job.test", "last_updated"),
		   resource.TestCheckResourceAttrSet("ansibleforms_job.test", "status"),
		),
      },
      // Delete testing automatically occurs in TestCase
    },
  })
}
