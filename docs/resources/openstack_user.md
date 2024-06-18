---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cleuracloud_openstack_user Resource - cleuracloud"
subcategory: ""
description: |-
  Creates a user in Cleura Cloud
---

# cleuracloud_openstack_user (Resource)

Creates a user in Cleura Cloud



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain_id` (String)
- `enabled` (Boolean)
- `name` (String)
- `projects` (Attributes List) (see [below for nested schema](#nestedatt--projects))

### Optional

- `default_project_id` (String)
- `description` (String)

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--projects"></a>
### Nested Schema for `projects`

Required:

- `id` (String)
- `roles` (List of String)