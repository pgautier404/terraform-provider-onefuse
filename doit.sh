#!/bin/bash

# destroy vm from previous run
echo 'setting environment...'
source ./setenv.sh

# destroys the prior vm and DELETE's the custom name 
echo 'running terraform destroy (for prior vm)'
terraform destroy

# clean, build, setup environment, init and apply the Terraform provider
./clean.sh
echo 'building go plugin...'
go build -o terraform-provider-sovlabs
echo 'running terraform init/plan/apply...'
terraform init
terraform plan
terraform apply
