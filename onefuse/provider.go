// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"scheme": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEFUSE_SCHEME", nil),
				Description: "OneFuse REST endpoint service http/https scheme",
			},
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEFUSE_ADDRESS", nil),
				Description: "OneFuse REST endpoint service host address",
			},
			"port": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEFUSE_PORT", nil),
				Description: "OneFuse REST endpoint service port number",
			},
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEFUSE_USER", nil),
				Description: "OneFuse REST endpoint user name",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ONEFUSE_PASSWORD", nil),
				Description: "OneFuse REST endpoint password",
			},
			"verify_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEFUSE_VERIFY_SSL", true),
				Description: "Verify SSL certificates for OneFuse endpoints",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"onefuse_naming":              resourceCustomNaming(),
			"onefuse_microsoft_ad_policy": resourceMicrosoftADPolicy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"onefuse_microsoft_endpoint": dataSourceMicrosoftEndpoint(),
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

func NewConfig(scheme string, address string, port string, user string, password string, verifySSL bool) *Config {
	return &Config{
		scheme:    scheme,
		address:   address,
		port:      port,
		user:      user,
		password:  password,
		verifySSL: verifySSL,
	}
}
