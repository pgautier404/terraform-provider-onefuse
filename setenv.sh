#!/bin/bash

export TF_LOG="TRACE"
export TF_LOG_PATH=terraform.log

export TF_VAR_sovlabs_address="localhost"
export TF_VAR_sovlabs_port=8080
export TF_VAR_sovlabs_user="root@sovlabs.local"
export TF_VAR_sovlabs_password="password"

export TF_VAR_vsphere_server=vcenter01.sovlabs.net

