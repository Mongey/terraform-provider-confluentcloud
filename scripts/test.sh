#!/bin/bash

set -ex

ARCH=$(go env GOARCH)
OS=$(go env GOOS)

make build

mv bin/${OS}-${ARCH}/terraform-provider-confluentcloud ~/.terraform.d/plugins/terraform-provider-confluentcloud
cd examples
terraform init
TF_LOG=debug terraform apply
