---
page_title: "confluentcloud_connector Resource - terraform-provider-confluentcloud"
subcategory: ""
description: |-
  
---

# Resource `confluentcloud_connector`





## Schema

### Required

- **cluster_id** (String) ID of containing cluster, e.g. lkc-abc123
- **config** (Map of String) Type-specific Configuration of connector. String keys and values
- **environment_id** (String) ID of containing environment, e.g. env-abc123
- **name** (String) The name of the connector

### Optional

- **config_sensitive** (Map of String) Sensitive part of connector configuration. String keys and values
- **id** (String) The ID of this resource.
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)


