#!/bin/bash

export TF_LOG="TRACE"
export TF_LOG_PATH=terraform.log

export TF_VAR_SOVLABS_ADDRESS="localhost"
export TF_VAR_SOVLABS_PORT=8080
export TF_VAR_SOVLABS_USER="root@sovlabs.local"
export TF_VAR_SOVLABS_PASSWORD="password"

export TF_VAR_VSPHERE_USER=<vsphere_user>
export TF_VAR_VSPHERE_PASSWORD=<password>
export TF_VAR_VSPHERE_SERVER=vcenter01.sovlabs.net

