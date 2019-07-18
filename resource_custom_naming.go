package main

import (
	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"log"
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
			"custom_name_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_suffix": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCustomNameCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("calling resourceCustomNameCreate")
	// call service to create/reserve custom name
	config := m.(Config)
	dnsSuffix := d.Get("dns_suffix").(string)
	cn := httpReserveCustomName(config, dnsSuffix)
	bindResource(d, cn)

	log.Println("custom name reserved: hostname=" + cn.Hostname + " dnsSuffix=" + cn.DnsSuffix)
	return resourceCustomNameRead(d, m)
}

func bindResource(d *schema.ResourceData, cn CustomName) {
	d.Set("custom_name_id", cn.Id)
	d.Set("hostname", cn.Hostname)
	d.Set("dns_suffix", cn.DnsSuffix)
}

type CustomName struct {
	Id        int
	Version   int
	Hostname  string
	DnsSuffix string
}

func httpReserveCustomName(config Config, dnsSuffix string) CustomName {
	refreshInputs := false
	address := config.address
	port := strconv.Itoa(config.port)
	url := "http://" + address + ":" + port + "/api/v1/customNames?refreshInputs=" + strconv.FormatBool(refreshInputs)
	log.Println("reserving custom name from " + url)

	postBody := "{\n\t\"namingStandard\": \"vraNamingStandard-{{Environment}}\",\n\t\"dnsSuffix\": \"" + dnsSuffix +
		"\",\n\t\"templateProperties\": {\n\t\t\"ownerName\": \"jsmith@company.com\",\n\t\t\"Environment\": \"dev\",\n\t\t\"OS\": \"Linux\",\n\t\t\"Application\": \"Web Servers\"\n\t},\n\t\"tenant\": \"defaultTenant\"\n}"
	payload := strings.NewReader(postBody)
	req, _ := http.NewRequest("POST", url, payload)
	log.Println("HTTP PAYLOAD to " + url + ":")
	log.Println(postBody)
	req.SetBasicAuth(config.user, config.password)
	setStandardHeaders(req)
	req.Header.Add("Host", address+":"+port)
	log.Println("REQUEST:")
	log.Println(req.Body)
	client := &http.Client{}
	res, err := client.Do(req);
	//res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	log.Println("HTTP RESULTS:")
	log.Println(body)
	var data CustomName
	json.Unmarshal(body, &data)
	res.Body.Close()
	return data
}

func setStandardHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	//req.Header.Add("Accept", "application/json")
	//req.Header.Add("Cache-Control", "no-cache")
	//req.Header.Add("accept-encoding", "gzip, deflate")
	//req.Header.Add("Connection", "keep-alive")
	//req.Header.Add("cache-control", "no-cache")
}

func httpGetCustomName(config Config, id int) CustomName {
	address := config.address
	port := strconv.Itoa(config.port)
	idString := strconv.Itoa(id)
	url := "http://" + address + ":" + port + "/api/v1/customNames/" + idString
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(config.user, config.password)
	setStandardHeaders(req)
	req.Header.Add("Host", address+":"+port)
	log.Println("REQUEST:")
	log.Println(req)
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
	config := m.(Config)
	id := d.Get("custom_name_id").(int)
	customName := httpGetCustomName(config, id)
	bindResource(d, customName)
	return nil
}

func resourceCustomNameUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceCustomNameDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
