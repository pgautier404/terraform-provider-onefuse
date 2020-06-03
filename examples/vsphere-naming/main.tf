
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

variable "vsphere_user" {
  type = string
}

variable "vsphere_password" {
  type = string
}

variable "vsphere_server" {
  type = string
}

provider "vsphere" {
  user = var.vsphere_user
  password = var.vsphere_password
  vsphere_server = var.vsphere_server

  # If you have a self-signed cert
  allow_unverified_ssl = true
  vim_keep_alive = 120
}

data "vsphere_datacenter" "dc" {
}

data "vsphere_datastore" "datastore" {
  name = "XtremIO_5t_datastore1"
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_resource_pool" "pool" {
  name = "Cluster1/Resources"
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_network" "network" {
  name = "dvs_SovLabs_331_10.30.31.0_24"
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_virtual_machine" "template" {
  name = "CentOS7"
  datacenter_id = data.vsphere_datacenter.dc.id
}

resource "vsphere_virtual_machine" "vm" {
  name = onefuse_naming.my-onefuse-name.name
  resource_pool_id = data.vsphere_resource_pool.pool.id
  datastore_id = data.vsphere_datastore.datastore.id

  num_cpus = 2
  memory = 1024
  guest_id = data.vsphere_virtual_machine.template.guest_id

  wait_for_guest_net_timeout = -1

  network_interface {
    network_id = data.vsphere_network.network.id
    adapter_type = data.vsphere_virtual_machine.template.network_interface_types[0]
  }

  disk {
    label = "disk0"
    size = data.vsphere_virtual_machine.template.disks.0.size
    eagerly_scrub = data.vsphere_virtual_machine.template.disks.0.eagerly_scrub
    thin_provisioned = data.vsphere_virtual_machine.template.disks.0.thin_provisioned
  }

  clone {
    template_uuid = data.vsphere_virtual_machine.template.id

    customize {
      linux_options {
        host_name = onefuse_naming.my-onefuse-name.name
        domain = onefuse_naming.my-onefuse-name.dns_suffix
      }

      network_interface {
        // windows requires per-network interface DNS settings so these may be ignored by linux
        dns_server_list = [
          "10.30.0.11",
          "10.30.0.12"]
        dns_domain = onefuse_naming.my-onefuse-name.dns_suffix
        ipv4_address = "10.30.31.203"
        ipv4_netmask = 24
      }

      ipv4_gateway = "10.30.31.1"
      dns_suffix_list = [
        onefuse_naming.my-onefuse-name.dns_suffix]
      // linux requires global DNS settings
      dns_server_list = [
        "10.30.0.11",
        "10.30.0.12"]
    }
  }

}
