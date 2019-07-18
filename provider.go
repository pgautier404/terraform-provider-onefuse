package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema {
			"address": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("SOVLABS_ADDRESS", nil),
				Description: "SovLabs REST endpoint service host address",
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("SOVLABS_PORT", nil),
				Description: "SovLabs REST endpoint service port number/integer",
			},
			"user": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("SOVLABS_USER", nil),
				Description: "SovLabs REST endpoint user name",
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("SOVLABS_PASSWORD", nil),
				Description: "SovLabs REST endpoint password",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sovlabs_custom_naming": resourceCustomNaming(),
		},
		ConfigureFunc: configureProvider,
	}
}

type Config struct {
	address string
	port int
	user string
	password string
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	return Config{
		address: d.Get("address").(string),
		port: d.Get("port").(int),
		user: d.Get("user").(string),
		password: d.Get("password").(string),
	}, nil
}