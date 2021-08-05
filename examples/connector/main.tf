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

resource "confluentcloud_connector" "connector" {
  name             = "pubsub-kafka-connector"
  environment_id   = "env-ab123"
  cluster_id       = "lkc-cd456"
  config           = {
    "name"                         = "pubsub-kafka-connector"
    "connector.class"              = "PubSubSource"
    "kafka.topic"                  = "kafka-topic1"
    "gcp.pubsub.project.id"        = "project-1234"
    "gcp.pubsub.subscription.id"   = "topic1-subscription1"
    "gcp.pubsub.topic.id"          = "topic1"
    "gcp.pubsub.max.retry.time"    = "5"
    "gcp.pubsub.message.max.count" = "1000"
    "errors.tolerance"             = "all"
    "tasks.max"                    = "1"
  }
  config_sensitive = {
    "kafka.api.key"               = <<kafka-api-key>>
    "kafka.api.secret"            = <<kafka-api-secret>>
    "gcp.pubsub.credentials.json" = <<gcp-service-account-key>
  }
}
