terraform {
  required_providers {
    ansibleforms = {
      source = "hashicorp.com/se/ansibleforms"
    }
  }
  required_version = ">= 0.0.1"
}

provider "ansibleforms" {
  username = "admin" # Publicly available by Ansible Forms
  password = "AnsibleForms!123" # Publicly available by Ansible Forms
  host = "https://127.0.0.1:8443" # Publicly available by Ansible Forms
}

resource "ansibleforms_job" "se" {
  form_name = "AXA Share Create"
  extravars = {
    region              = "myregion"
    opco                = "myopco"
    svm_name            = "mysvm_name"
    exposure            = "myexposure"
    env                 = "myenv"
    dataclass           = "mydataclass"
    share_name          = "myshare_name"
    accountid           = "myaccountid"
    size                = "mysize"
    protection_required = "myprotection_required"
  }
  credentials = {
    ontap_cred = "myontap_cred"
    bind_cred  = "mybind_cred"
  }
}

output "ansibleforms_job" {
  value = ansibleforms_job.se
}