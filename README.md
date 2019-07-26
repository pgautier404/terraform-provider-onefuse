
# Terraform Provider for SovLabs "Defender"

Terraform Provider to access SovLabs "Defender" platform


## Using the Provider

In your .tf file, declare the "sovlabs" provider:

```
provider "sovlabs" {
  address = "my-sovlabs-service-address.mycompany.com"
  port = 8080
  user = "my-user"
  password = "my-password"
}
```

(The clear-text password is used here only for demonstration purposes,
see Terraform documentation for ways to secure secrets)

Dynamic properties to be used in the naming standard template may be configured
as a Terraform map type as follows:
```
variable "template_properties_map" {
  type = "map"
  default = {
    // these are dynamic properties,
    // they can be anything that is defined in the Liquid template
    ownerName = "jsmith@company.com"
    environment = "dev"
    os = "Linux"
    application = "Web Servers"
  }
}
```

Declare the custom-naming resource, note how the built-in function
"jsonencode" is used to encode the map properties as a string:
```
resource "sovlabs_custom_naming" "my-custom-name" {
  template_properties = jsonencode(var.template_properties_map)
  dns_suffix = "bluecat90.sovlabs.net"
  hostname = ""
}
```

The "hostname" parameter must be configured as empty string, the value
will be retrieved by the provider.

This example shows how template properties may also be loaded from a stand-alone .json file:
```
resource "sovlabs_custom_naming" "my-custom-name" {
  template_properties = file("template_properties.json")
  dns_suffix = "bluecat90.sovlabs.net"
  hostname = ""
}
```

The included example main.tf file also shows how the VMWare VSphere provider may
be used to provision the custom name acquired from the SovLabs provider.  See the 
terraform-provider-vsphere documentation for more information.

