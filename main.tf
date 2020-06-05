
variable "onefuse_user" {
  type = string
}

variable "onefuse_password" {
  type = string
}

variable "onefuse_address" {
  type = string
}

variable "onefuse_port" {
  type = string
}

provider "onefuse" {
  address     = var.onefuse_address
  port        = var.onefuse_port
  user        = var.onefuse_user
  password    = var.onefuse_password
  verify_ssl  = false
}

variable "onefuse_naming_policy_id" {
  type = string
}

variable "onefuse_dns_suffix" {
  type = string
}

resource "onefuse_naming" "my-onefuse-name" {
  naming_policy_id        = var.onefuse_naming_policy_id
  dns_suffix              = var.onefuse_dns_suffix
  template_properties     = {
      "ownerName"               = "jsmith@company.com"
      "Environment"             = "dev"
      "OS"                      = "Linux"
      "Application"             = "Web Servers"
      "suffix"                  = "sovlabs.net"
      "tenant"                  =  "mytenant"
  }
}