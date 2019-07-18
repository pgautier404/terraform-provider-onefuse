package main

import (
	"log"
	"strconv"
	"testing"
)


func TestHttpReserveCustomName(t *testing.T) {
	config := Config {
		address: "127.0.0.1",
		port: 8080,
		user: "root@defender.local",
		password: "password",
	}
	cn := httpReserveCustomName(config, "sovlabs.net")
	log.Println(strconv.Itoa(cn.Id) + ": " + cn.Hostname + "." + cn.DnsSuffix + " version:" + strconv.Itoa(cn.Version))
}