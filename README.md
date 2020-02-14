# `terraform-plugin-confluent-cloud`

A [Terraform][1] plugin for managing [Confluent Cloud Kafka Clusters][2].

## Installation

Download and extract the [latest release](https://github.com/Mongey/terraform-provider-confluent-cloud/releases/latest) to
your [terraform plugin directory][third-party-plugins] (typically `~/.terraform.d/plugins/`)

## Example

Configure the provider directly, or set the ENV variables `CONFLUENT_CLOUD_USERNAME` &`CONFLUENT_CLOUD_PASSWORD`

```hcl
provider "confluentcloud" {}

resource "confluentcloud_environment" "test" {
  name = "provider-test"
}

resource "confluentcloud_kafka_cluster" "test" {
  name             = "provider-test"
  environment_id   = confluentcloud_environment.test.id
  service_provider = "aws"
  region           = "eu-west-1"
  availability     = "LOW"
}

resource "confluentcloud_api_key" "provider_test" {
  cluster_id = confluentcloud_kafka_cluster.test.id
  environment_id   = confluentcloud_environment.test.id
}

provider "kafka" {
  bootstrap_servers = [replace(confluentcloud_kafka_cluster.test.bootstrap_servers, "SASL_SSL://", "")]

  tls_enabled    = true
  sasl_username  = confluentcloud_api_key.provider_test.key
  sasl_password  = confluentcloud_api_key.provider_test.secret
  sasl_mechanism = "plain"
}

resource "kafka_topic" "syslog" {
  name               = "syslog2"
  replication_factor = 3
  partitions         = 1
}

output "kafka_url" {
  value = replace(confluentcloud_kafka_cluster.test.bootstrap_servers, "SASL_SSL://", "")
}

output "key" {
  value = confluentcloud_api_key.provider_test.key
}

output "secret" {
  value = confluentcloud_api_key.provider_test.secret
}
```

[1]: https://www.terraform.io
[2]: https://confluent.cloud
[third-party-plugins]: https://www.terraform.io/docs/configuration/providers.html#third-party-plugins
