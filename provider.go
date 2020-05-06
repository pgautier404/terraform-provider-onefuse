package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"scheme": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("FUSE_SCHEME", "http"),
				Description: "Fuse REST endpoint scheme (http or https)",
			},
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FUSE_ADDRESS", nil),
				Description: "Fuse REST endpoint service host address",
			},
			"port": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FUSE_PORT", nil),
				Description: "Fuse REST endpoint service port number",
			},
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FUSE_USER", nil),
				Description: "Fuse REST endpoint user name",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("FUSE_PASSWORD", nil),
				Description: "Fuse REST endpoint password",
			},
			"verify_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("FUSE_VERIFY_SSL", true),
				Description: "Verify SSL certificates for Fuse endpoints",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"fuse_naming": resourceCustomNaming(),
		},
		ConfigureFunc: configureProvider,
	}
}

type Config struct {
	scheme    string
	address   string
	port      string
	user      string
	password  string
	verifySSL bool
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	return Config{
		scheme:    d.Get("scheme").(string),
		address:   d.Get("address").(string),
		port:      d.Get("port").(string),
		user:      d.Get("user").(string),
		password:  d.Get("password").(string),
		verifySSL: d.Get("verify_ssl").(bool),
	}, nil
}
