#!/bin/bash

set -ex

ARCH=go env GOARCH
OS=go env GOOS

go build

mv bin/${OS}-${ARCH}/terraform-provider-confluentcloud ~/.terraform.d/plugins/${OS}_${ARCH}/
cd examples
terraform init
terraform plan
terraform output
