package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func resourceCustomNaming() *schema.Resource {
	return &schema.Resource{
		Create: resourceCustomNameCreate,
		Read:   resourceCustomNameRead,
		Update: resourceCustomNameUpdate,
		Delete: resourceCustomNameDelete,
		Schema: map[string]*schema.Schema{
			"refresh_inputs": {
				Type:     schema.TypeBool,
				Default:  false,
				Required: false,
				Optional: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_suffix": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCustomNameCreate(d *schema.ResourceData, m interface{}) error {
	fmt.Println("calling resourceCustomNameCreate")
	// call service to create/reserve custom name

	cn := httpReserveCustomName(d, m)
	d.SetId(cn.Id) // numeric id assigned by back-end
	d.Set("hostname", cn.Hostname)
	d.Set("dns_suffix", cn.DnsSuffix)

	fmt.Println("custom name reserved: hostname=" + cn.Hostname + " dnsSuffix=" + cn.DnsSuffix)
	return resourceCustomNameRead(d, m)
}

type CustomName struct {
	Id        string
	Version   int
	Hostname  string
	DnsSuffix string
}

func httpReserveCustomName(d *schema.ResourceData, m interface{}) CustomName {
	config := m.(Config)
	refreshInputs := d.Get("refresh_inputs").(bool)
	address := config.address
	port := strconv.Itoa(config.port)
	url := "http://" + address + ":" + port + "/api/v1/customNames?refreshInputs=" + strconv.FormatBool(refreshInputs)

	fmt.Println("reserving custom name from " + url)

	payload := strings.NewReader("{\n\t\"namingStandard\": {\n\t\t\"id\": 1,\n\t\t\"version\": 0,\n\t\t\"namingSequences\": [\n\t\t  \t{\n\t\t\t  \"format\": \"vra73cls/##/\",\n\t\t\t  \"initialvalue\": \"vra73cls01\",\n\t\t\t  \"reuse\": false,\n\t\t\t  \"name\": \"vraSequence\",\n\t\t\t  \"tenant\": {\n\t\t\t  \t \"id\": 1,\n\t\t\t  \t \"version\": 0,\n\t\t\t  \t \"name\": \"defaultTenant\"\n\t\t\t  },\n\t\t\t  \"version\": 0,\n\t\t\t  \"type\": \"pattern\",\n\t\t\t  \"uniquekey\": \"{{ ownerName | split: \\\"@\\\" | first |substring: 0,2 | downcase}}{{ Environment | substring: 0,1| downcase }}{{ OS | substring: 0,1 | downcase }}{{ Application | substring: 0,3 | downcase }}\"\n\t\t\t}\n\t  \t],\n\t  \"name\": \"vraNamingStandard-{{Environment}}\",\n\t  \"template\": \"{{sequence.vraSequence}}\",\n\t  \"tenant\": {\n\t        \"id\": 1,\n\t        \"version\": 0,\n\t        \"name\": \"defaultTenant\"\n\t    }\n\t},\n\t\"dnsSuffix\": \"sovlabs.net\",\n\t\"templateProperties\": {\n\t\t\"ownerName\": \"jsmith@company.com\",\n        \"Environment\": \"dev\",\n        \"OS\": \"Linux\",\n        \"Application\": \"Web Servers\"\n\t},\n\t\"tenant\": {\n\t  \t \"id\": 1,\n\t  \t \"version\": 0,\n\t  \t \"name\": \"defaultTenant\"\n\t}\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.SetBasicAuth(config.user, config.password)

	setStandardHeaders(req)
	req.Header.Add("Host", address+":"+port)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var data CustomName
	json.Unmarshal(body, &data)
	res.Body.Close()
	return data
}

func setStandardHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("accept-encoding", "gzip, deflate")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")
}


func httpGetCustomName(d *schema.ResourceData, m interface{}) CustomName {
	config := m.(Config)
	address := config.address
	port := strconv.Itoa(config.port)
	id := d.Id()
	url := "http://" + address + ":" + port + "/api/v1/customNames/" + id

	req, _ := http.NewRequest("GET", url, nil)

	req.SetBasicAuth(config.user, config.password)

	setStandardHeaders(req)
	req.Header.Add("Host", address+":"+port)

	res, _ := http.DefaultClient.Do(req)
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err.Error())
	}
	var data CustomName
	json.Unmarshal(body, &data)
	res.Body.Close()
	return data

}

func resourceCustomNameRead(d *schema.ResourceData, m interface{}) error {
	customName := httpGetCustomName(d, m)
	d.SetId(customName.Id)
	d.Set("hostname", customName.Hostname)
	d.Set("dns_suffix", customName.DnsSuffix)
	return nil
}

func resourceCustomNameUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceCustomNameRead(d, m)
}

func resourceCustomNameDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
