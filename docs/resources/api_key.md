---
page_title: "confluentcloud_api_key Resource - terraform-provider-confluentcloud"
subcategory: ""
description: |-
  
---

# Resource `confluentcloud_api_key`





## Schema

### Required

- **environment_id** (String) Environment ID

### Optional

- **cluster_id** (String)
- **description** (String) Description
- **id** (String) The ID of this resource.
- **logical_clusters** (List of String) Logical Cluster ID List to create API Key
- **user_id** (Number) User ID

### Read-only

- **key** (String)
- **secret** (String, Sensitive)


