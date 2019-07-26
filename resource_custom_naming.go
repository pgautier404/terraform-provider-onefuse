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
			"template_properties" : {
				Type: schema.TypeString,
				Optional: true,
			},
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

	// debug the resource(tf state) and interface(config) data
	// requires 3rd party library "spew": import "github.com/davecgh/go-spew/spew"
	//log.Println("calling resourceCustomNameCreate")
	//log.Println("RESOURCE DATA:")
	//log.Println(spew.Sprint(d))
	//log.Print("END RESOURCE DATA")
	//log.Println("INTERFACE DATA:")
	//log.Println(spew.Sprint(m))
	//log.Print("END INTERFACE DATA")

	// call service to create/reserve custom name
	config := m.(Config)
	dnsSuffix := d.Get("dns_suffix").(string)

	// if we want to get the properties as a map directly from .tf without json encoding
	// then we can do this:
	// templateProperties := d.Get("template_properties").(map[string]interface{})
	// jsonBytes, err := json.Marshal(templateProperties)
	// jsonString := string(jsonBytes)

	// get template_properties as string
	jsonString := d.Get("template_properties").(string)

	cn, err := httpReserveCustomName(config, jsonString, dnsSuffix)
	if err != nil {
		return errors.WithMessage(err, "Failed to reseve custom name")
	}
	if err := bindResource(d, cn); err != nil {
		return errors.WithMessage(err, "cannot bind custom name reservation to resource, perhaps the API has changed?")
	}
	return nil
}

func bindResource(d *schema.ResourceData, cn CustomName) error {

	// setting the ID is REALLY necessary here
	// we use the full host.domain instead of a numeric ID as it will remain a consistent composite key
	d.SetId(cn.Name + "." + cn.DnsSuffix)

	if err := d.Set("custom_name_id", cn.Id); err != nil {
		return errors.WithMessage(err, "cannot set custom_name_id")
	}
	if err := d.Set("hostname", cn.Name); err != nil {
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
	Name      string
	DnsSuffix string
}

func httpReserveCustomName(config Config, templateProperties string, dnsSuffix string) (CustomName, error) {
	address := config.address
	port := strconv.Itoa(config.port)
	url := "http://" + address + ":" + port + "/api/v1/customNames?refreshInputs=" + strconv.FormatBool(false)
	log.Println("reserving custom name from " + url + "  dnsSuffix=" + dnsSuffix)
	postBody := "{\n\t\"namingStandard\": \"vraNamingStandard-{{Environment}}\",\n\t\"dnsSuffix\": \"" + dnsSuffix +
		"\",\n\t\"templateProperties\": " + templateProperties +
		//"{\n\t\t\"ownerName\": \"jsmith@company.com\",\n\t\t\"Environment\": \"dev\",\n\t\t\"OS\": \"Linux\",\n\t\t\"Application\": \"Web Servers\"\n\t}" +
		",\n\t\"tenant\": \"defaultTenant\"\n}"
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
	if err := json.Unmarshal(body, &data); err != nil {
		return data, errors.WithMessage(err, "Failed to unmarshal HTTP POST JSON response")
	}
	if err := res.Body.Close(); err != nil {
		return data, errors.WithMessage(err, "Failed to close HTTP response body stream")
	}
	log.Println("custom name reserved: " +
		"custom_name_id=" + strconv.Itoa(data.Id) +
		" hostname=" + data.Name +
		" dnsSuffix=" + data.DnsSuffix)
	return data, nil
}

func setStandardHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("accept-encoding", "gzip, deflate")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")
}

func httpGetCustomName(config Config, id int) (CustomName, error) {
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
	var data CustomName
	if err != nil {
		return data, errors.WithMessage(err, "Failed to unmarshal response JSON")
	}
	log.Println("HTTP GET RESULTS:")
	log.Println(string(body))
	if err = json.Unmarshal(body, &data); err != nil {
		return data, errors.WithMessage(err, "Failed to unmarshal response JSON")
	}
	if err := res.Body.Close(); err != nil {
		return data, errors.WithMessage(err, "Failed to close response body stream")
	}
	return data, nil
}

func resourceCustomNameRead(d *schema.ResourceData, m interface{}) error {
	config := m.(Config)
	id := d.Get("custom_name_id").(int)
	customName, err := httpGetCustomName(config, id)
	if err != nil {
		return errors.WithMessage(err, "Failed to get custom name from REST API")
	}
	if err := bindResource(d, customName); err != nil {
		return errors.WithMessage(err, "Failed to bind REST API response to terraform resource, perhaps API has changed?")
	}
	return nil
}

func resourceCustomNameUpdate(_ *schema.ResourceData, _ interface{}) error {
	return nil
}

func resourceCustomNameDelete(_ *schema.ResourceData, _ interface{}) error {
	return nil
}
