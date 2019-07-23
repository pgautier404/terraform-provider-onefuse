#!/bin/bash

export TF_LOG="TRACE"
export TF_LOG_PATH=/tmp/terraform-log
export SOVLABS_ADDRESS="localhost"
export SOVLABS_PORT=8080
export SOVLABS_USER="root@sovlabs.local"
export SOVLABS_PASSWORD="password"

export TF_VAR_vsphere_user=jhitt
export TF_VAR_vsphere_password=Disco0915
export TF_VAR_vsphere_server=vcenter01.sovlabs.net

