#!/bin/bash

# clean, build, setup environment, init and apply the Terraform provider

rm -f terraform.tfstate
echo 'building go plugin...'
go build -o terraform-provider-sovlabs
echo 'setting environment...'
source ./setenv.sh
echo 'running terraform...'
terraform init
terraform apply

