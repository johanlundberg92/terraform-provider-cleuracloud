---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cleuracloud_ccp_user Resource - cleuracloud"
subcategory: ""
description: |-
  Creates a CCP user in Cleura Cloud
---

# cleuracloud_ccp_user (Resource)

Creates a CCP user in Cleura Cloud



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `email` (String)
- `name` (String)

### Optional

- `first_name` (String)
- `last_name` (String)
- `privileges` (Attributes) (see [below for nested schema](#nestedatt--privileges))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--privileges"></a>
### Nested Schema for `privileges`

Optional:

- `openstack` (Attributes) (see [below for nested schema](#nestedatt--privileges--openstack))
- `users` (Attributes) (see [below for nested schema](#nestedatt--privileges--users))

<a id="nestedatt--privileges--openstack"></a>
### Nested Schema for `privileges.openstack`

Required:

- `type` (String)


<a id="nestedatt--privileges--users"></a>
### Nested Schema for `privileges.users`

Required:

- `type` (String)
