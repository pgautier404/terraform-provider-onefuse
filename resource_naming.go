package main

import (
	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
	"strconv"
	"strings"
)

func resourceCustomName() *schema.Resource {
	return &schema.Resource{
		Create: resourceCustomNameCreate,
		Read:   resourceCustomNameRead,
		Update: resourceCustomNameUpdate,
		Delete: resourceCustomNameDelete,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"user": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"refreshInputs": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Required: false,
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"dnsSuffix": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
		},
	}
}

func resourceCustomNameCreate(d *schema.ResourceData, m interface{}) error {
	// call service to create/reserve custom name
	cn := httpReserveCustomName(d)
	// set resource id to the full hostname
	d.SetId(cn.Id)
	d.Set("hostname", cn.Hostname)
	d.Set("dnsSuffix", cn.DnsSuffix)
	return resourceCustomNameRead(d, m)
}

type custom_name struct {
	Id        string
	Version   int
	Hostname  string
	DnsSuffix string
}

func httpReserveCustomName(d *schema.ResourceData) (custom_name) {
	hostname := d.Get("hostname").(string)
	port := strconv.Itoa(d.Get("port").(int))
	refreshInputs := strconv.FormatBool(d.Get("refreshInputs").(bool))
	url := "http://" + hostname + ":" + port + "/api/v1/customNames?refreshInputs=" + refreshInputs

	payload := strings.NewReader("{\n\t\"namingStandard\": {\n\t\t\"id\": 1,\n\t\t\"version\": 0,\n\t\t\"namingSequences\": [\n\t\t  \t{\n\t\t\t  \"format\": \"vra73cls/##/\",\n\t\t\t  \"initialvalue\": \"vra73cls01\",\n\t\t\t  \"reuse\": false,\n\t\t\t  \"name\": \"vraSequence\",\n\t\t\t  \"tenant\": {\n\t\t\t  \t \"id\": 1,\n\t\t\t  \t \"version\": 0,\n\t\t\t  \t \"name\": \"defaultTenant\"\n\t\t\t  },\n\t\t\t  \"version\": 0,\n\t\t\t  \"type\": \"pattern\",\n\t\t\t  \"uniquekey\": \"{{ ownerName | split: \\\"@\\\" | first |substring: 0,2 | downcase}}{{ Environment | substring: 0,1| downcase }}{{ OS | substring: 0,1 | downcase }}{{ Application | substring: 0,3 | downcase }}\"\n\t\t\t}\n\t  \t],\n\t  \"name\": \"vraNamingStandard-{{Environment}}\",\n\t  \"template\": \"{{sequence.vraSequence}}\",\n\t  \"tenant\": {\n\t        \"id\": 1,\n\t        \"version\": 0,\n\t        \"name\": \"defaultTenant\"\n\t    }\n\t},\n\t\"dnsSuffix\": \"sovlabs.net\",\n\t\"templateProperties\": {\n\t\t\"ownerName\": \"jsmith@company.com\",\n        \"Environment\": \"dev\",\n        \"OS\": \"Linux\",\n        \"Application\": \"Web Servers\"\n\t},\n\t\"tenant\": {\n\t  \t \"id\": 1,\n\t  \t \"version\": 0,\n\t  \t \"name\": \"defaultTenant\"\n\t}\n}")

	req, _ := http.NewRequest("POST", url, payload)
	user := d.Get("user").(string)
	password := d.Get("password").(string)
	req.SetBasicAuth(user, password)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", hostname + ":" + port)
	req.Header.Add("accept-encoding", "gzip, deflate")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	var data custom_name
	decoder.Decode(&data)
	return data
}

func httpGetCustomName(d *schema.ResourceData) (custom_name) {
	hostname := d.Get("hostname").(string)
	port := strconv.Itoa(d.Get("port").(int))
	id := d.Id()
	url := "http://" + hostname + ":" + port + "/api/v1/customNames/" + id

	req, _ := http.NewRequest("GET", url, nil)

	user := d.Get("user").(string)
	password := d.Get("password").(string)
	req.SetBasicAuth(user, password)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", hostname + ":" + port)
	req.Header.Add("accept-encoding", "gzip, deflate")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	var data custom_name
	decoder.Decode(&data)
	return data

}

func resourceCustomNameRead(d *schema.ResourceData, m interface{}) error {
	custom_name := httpGetCustomName(d)
	d.Set("hostname", custom_name.Hostname)
	d.Set("dnsSuffix", custom_name.DnsSuffix)
	return nil
}

func resourceCustomNameUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceCustomNameRead(d, m)
}

func resourceCustomNameDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
