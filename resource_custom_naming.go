package main

import (
	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
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
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
	return nil
}

func bindResource(d *schema.ResourceData, cn CustomName) error {

	// setting the ID is REALLY necessary here
	// we use the FQDN instead of the numeric ID as it is more likely to remain consistent as a composite key in TF
	d.SetId(cn.Hostname + "." + cn.DnsSuffix)

	if err := d.Set("custom_name_id", cn.Id); err != nil {
		return errors.WithMessage(err, "cannot set custom_name_id")
	}
	if err := d.Set("hostname", cn.Hostname); err != nil {
		return errors.WithMessage(err, "cannot set hostname")
	}
	if err := d.Set("dns_suffix", cn.DnsSuffix); err != nil {
		return errors.WithMessage(err, "cannot set dns_suffix")
	}
	return nil
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
	log.Println("reserving custom name from " + url + "  dnsSuffix=" + dnsSuffix)
	postBody := "{\n\t\"namingStandard\": \"vraNamingStandard-{{Environment}}\",\n\t\"dnsSuffix\": \"" + dnsSuffix +
		"\",\n\t\"templateProperties\": {\n\t\t\"ownerName\": \"jsmith@company.com\",\n\t\t\"Environment\": \"dev\",\n\t\t\"OS\": \"Linux\",\n\t\t\"Application\": \"Web Servers\"\n\t},\n\t\"tenant\": \"defaultTenant\"\n}"
	payload := strings.NewReader(postBody)
	log.Println("CONFIG:")
	log.Println(config)
	req, _ := http.NewRequest("POST", url, payload)
	log.Println("HTTP PAYLOAD to " + url + ":")
	log.Println(postBody)
	req.SetBasicAuth(config.user, config.password)
	setStandardHeaders(req)
	req.Header.Add("Host", address+":"+port)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	log.Println("HTTP POST RESULTS:")
	log.Println(string(body))
	var data CustomName
	json.Unmarshal(body, &data)
	res.Body.Close()
	log.Println("custom name reserved: " +
		"custom_name_id=" + strconv.Itoa(data.Id) +
		" hostname=" + data.Hostname +
		" dnsSuffix=" + data.DnsSuffix)
	return data
}

func setStandardHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("accept-encoding", "gzip, deflate")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")
}

func httpGetCustomName(config Config, id int) CustomName {
	address := config.address
	port := strconv.Itoa(config.port)
	idString := strconv.Itoa(id)
	url := "http://" + address + ":" + port + "/api/v1/customNames/" + idString
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
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
	log.Println("HTTP GET RESULTS:")
	log.Println(string(body))
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
