---
page_title: "confluentcloud_kafka_cluster Resource - terraform-provider-confluentcloud"
subcategory: ""
description: |-
  
---

# Resource `confluentcloud_kafka_cluster`





## Schema

### Required

- **availability** (String) LOW(single-zone) or HIGH(multi-zone)
- **environment_id** (String) Environment ID
- **name** (String) The name of the cluster
- **region** (String) where
- **service_provider** (String) AWS / GCP

### Optional

- **cku** (Number) cku
- **deployment** (Map of String) Deployment settings.  Currently only `sku` is supported.
- **id** (String) The ID of this resource.
- **network_egress** (Number) Network egress limit(MBps)
- **network_ingress** (Number) Network ingress limit(MBps)
- **storage** (Number) Storage limit(GB)

### Read-only

- **bootstrap_servers** (String)


