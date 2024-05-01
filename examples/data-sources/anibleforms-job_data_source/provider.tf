terraform {
  required_providers {
    ansibleforms = {
      source = "hashicorp.com/se/ansibleforms"
    }
  }
  required_version = ">= 0.0.1"
}

provider "ansibleforms" {
  connection_profiles = [
    {
      name = "cluster1"
      hostname = "********219"
      username = var.username
      password = var.password
      hostname = "127.0.0.1:8443" # Publicly available by Ansible Forms
      validate_certs = var.validate_certs
    }
  ]
}
