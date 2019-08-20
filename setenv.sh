#!/bin/bash

export TF_LOG="TRACE"
export TF_LOG_PATH=terraform.log

export TF_VAR_sovlabs_address="nightly01.sovlabs.net"
export TF_VAR_sovlabs_port=3033
export TF_VAR_sovlabs_user="mmaxwell@sovlabs.local"
export TF_VAR_sovlabs_password="password"

export TF_VAR_vsphere_server=vcenter01.sovlabs.net
export TF_VAR_vsphere_user=vrasvc
export TF_VAR_vsphere_password=VmwareS0v
