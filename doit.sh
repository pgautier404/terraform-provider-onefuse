#!/bin/bash

# destroy vm from previous run
echo 'setting environment...'
source ./setenv.sh
echo 'running terraform destroy (for prior vm)'
terraform destroy

# clean, build, setup environment, init and apply the Terraform provider
rm -f terraform.tfstate*
rm -rf /tmp/tf-state*
rm /tmp/terraform-log
echo 'building go plugin...'
go build -o terraform-provider-sovlabs
echo 'running terraform init/plan/apply...'
terraform init
terraform plan
terraform apply

