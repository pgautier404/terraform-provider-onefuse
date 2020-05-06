# terraform-provider-fuse
Terraform Provider for integrating with SovLabs Fuse.

## Sample Terraform Configuration
To get started with the Terraform Provider for SovLabs Fuse, put the following into a file called `main.tf`.

Fill in the `provider "fuse"` section with details about your SovLabs Fuse instance.

```hcl
provider "fuse" {
  address     = "localhost"
  port        = "8000"
  user        = "admin"
  password    = "my-password"
  scheme      = "http"
  verify_ssl  = false
}

resource "fuse_naming" "my-fuse-name" {
  naming_policy_id        = "2"
  dns_suffix              = "sovlabs.net"
  workspace_id            = "6"
  template_properties     = {
      "ownerName"               = "jsmith@company.com"
      "Environment"             = "dev"
      "OS"                      = "Linux"
      "Application"             = "Web Servers"
      "suffix"                  = "sovlabs.net"
      "tenant"                  =  "mytenant"
  }
}
```