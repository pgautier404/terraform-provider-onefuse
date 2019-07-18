package main

import (
	"log"
	"strconv"
	"testing"
)

func TestHttpReserveCustomName(t *testing.T) {
	config := GetConfig()
	cn := httpReserveCustomName(config, "sovlabs.net")
	log.Println(strconv.Itoa(cn.Id) + ": " + cn.Hostname + "." + cn.DnsSuffix + " version:" + strconv.Itoa(cn.Version))
	if cn.Id <= 0 {
		t.Errorf("customName.Id=%d; want > 0", cn.Id)
	}
	if cn.DnsSuffix != "sovlabs.net" {
		t.Errorf("customName.DnsSuffix=%d; want sovlabs.net", cn.DnsSuffix)
	}
	if cn.Hostname == "" {
		t.Errorf("customName.Hostname=%d; want non-empty string", cn.Hostname)
	}
	if cn.Version <= 0 {
		t.Errorf("customName.Version=%d; want > 0", cn.Version)
	}
}

func TestHttpGetCustomName(t *testing.T) {
	config := GetConfig()
	cn1 := httpReserveCustomName(config, "sovlabs.net")
	cn2 := httpGetCustomName(config, cn1.Id)
	if cn1.Id != cn2.Id {
		t.Error("Reserved customName.Id does not match after retrieval")
	}
	if cn1.Hostname != cn2.Hostname {
		t.Error("Reserved customName.Hostname does not match after retrieval")
	}
	if cn1.DnsSuffix != cn2.DnsSuffix {
		t.Error("Reserved customName.DnsSuffix does not match after retrieval")
	}
	if cn1.Version != cn2.Version {
		t.Error("Reserved customName.Version does not match after retrieval")
	}

}

func GetConfig() Config {
	config := Config{
		address:  "localhost",
		port:     8080,
		user:     "root@sovlabs.local",
		password: "password",
	}
	return config
}
