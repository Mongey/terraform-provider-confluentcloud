terraform {
  required_providers {
    kafka = {
      source  = "Mongey/kafka"
      version = "0.2.11"
    }
    confluentcloud = {
      source = "Mongey/confluentcloud"
    }
  }
}

provider "confluentcloud" {}

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

  depends_on = [confluentcloud_kafka_cluster.test]
}

# api key for a kafka cluster
resource "confluentcloud_api_key" "provider_test" {
  cluster_id     = confluentcloud_kafka_cluster.test.id
  environment_id = confluentcloud_environment.environment.id
}

# api key for schema registry
resource "confluentcloud_api_key" "schema-registry" {
  logical_clusters = [
    confluentcloud_schema_registry.test.id
  ]
  environment_id = confluentcloud_environment.environment.id
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
