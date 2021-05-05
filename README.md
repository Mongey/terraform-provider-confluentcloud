# `terraform-plugin-confluentcloud`

A [Terraform][1] plugin for managing [Confluent Cloud Kafka Clusters][2].

## Installation

Download and extract the [latest release](https://github.com/Mongey/terraform-provider-confluentcloud/releases/latest) to
your [terraform plugin directory][third-party-plugins] (typically `~/.terraform.d/plugins/`) or define the plugin in the required_providers block.

```hcl
terraform {
  required_providers {
    confluentcloud = {
      source = "Mongey/confluentcloud"
    }
  }
}
```

## Example

Configure the provider directly, or set the ENV variables `CONFLUENT_CLOUD_USERNAME` &`CONFLUENT_CLOUD_PASSWORD`

```hcl
terraform {
  required_providers {
    confluentcloud = {
      source = "Mongey/confluentcloud"
    }
    kafka = {
      source  = "Mongey/kafka"
      version = "0.2.11"
    }
  }
}

provider "confluentcloud" {
  username = "ccloud@example.org"
  password = "hunter2"
}

resource "confluentcloud_environment" "environment" {
  name = "production"
}

resource "confluentcloud_kafka_cluster" "test" {
  name             = "provider-test"
  service_provider = "aws"
  region           = "eu-west-1"
  availability     = "LOW"
  environment_id   = confluentcloud_environment.environment.id
  deployment = {
    sku = "BASIC"
  }
  network_egress  = 100
  network_ingress = 100
  storage         = 5000
}

resource "confluentcloud_schema_registry" "test" {
  environment_id   = confluentcloud_environment.environment.id
  service_provider = "aws"
  region           = "EU"

  # Requires at least one kafka cluster to enable the schema registry in the environment.
  depends_on = [confluentcloud_kafka_cluster.test]
}

resource "confluentcloud_api_key" "provider_test" {
  cluster_id     = confluentcloud_kafka_cluster.test.id
  environment_id = confluentcloud_environment.environment.id
}

resource "confluentcloud_service_account" "test" {
  name           = "test"
  description    = "service account test"
}

locals {
  bootstrap_servers = [replace(confluentcloud_kafka_cluster.test.bootstrap_servers, "SASL_SSL://", "")]
}

provider "kafka" {
  bootstrap_servers = local.bootstrap_servers

  tls_enabled    = true
  sasl_username  = confluentcloud_api_key.provider_test.key
  sasl_password  = confluentcloud_api_key.provider_test.secret
  sasl_mechanism = "plain"
  timeout        = 10
}

resource "kafka_topic" "syslog" {
  name               = "syslog"
  replication_factor = 3
  partitions         = 1
  config = {
    "cleanup.policy" = "delete"
  }
}

output "kafka_url" {
  value = local.bootstrap_servers
}

output "key" {
  value     = confluentcloud_api_key.provider_test.key
  sensitive = true
}

output "secret" {
  value     = confluentcloud_api_key.provider_test.secret
  sensitive = true
}
```

## Importing existing resources

This provider supports importing existing Confluent Cloud resources via [`terraform import`][3].

Most resource types use the import IDs returned by the [`ccloud`][4] CLI.
`confluentcloud_kafka_cluster` and `confluentcloud_schema_registry` can be imported using `<environment ID>/<cluster ID>`.

[1]: https://www.terraform.io
[2]: https://confluent.cloud
[3]: https://www.terraform.io/docs/cli/import/index.html
[4]: https://docs.confluent.io/ccloud-cli/current/index.html
[third-party-plugins]: https://www.terraform.io/docs/configuration/providers.html#third-party-plugins
