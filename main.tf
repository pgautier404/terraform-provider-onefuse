
variable "fuse_user" {
  type = string
}

variable "fuse_password" {
  type = string
}

variable "fuse_address" {
  type = string
}

variable "fuse_port" {
  type = string
}

variable "fuse_scheme" {
  type = string
}

provider "fuse" {
  address     = var.fuse_address
  port        = var.fuse_port
  user        = var.fuse_user
  password    = var.fuse_password
  scheme      = var.fuse_scheme
  verify_ssl  = false
}

variable "fuse_naming_policy_id" {
  type = string
}

variable "fuse_dns_suffix" {
  type = string
}

resource "fuse_naming" "my-fuse-name" {
  naming_policy_id        = var.fuse_naming_policy_id
  dns_suffix              = var.fuse_dns_suffix
  template_properties     = {
      "ownerName"               = "jsmith@company.com"
      "Environment"             = "dev"
      "OS"                      = "Linux"
      "Application"             = "Web Servers"
      "suffix"                  = "sovlabs.net"
      "tenant"                  =  "mytenant"
  }
}