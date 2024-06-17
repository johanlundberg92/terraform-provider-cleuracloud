---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cleura Provider"
subcategory: ""
description: |-
  Interact with Cleura.
---

# cleura Provider

Interact with Cleura.

## Example Usage

```terraform
# Configuration-based authentication
provider "hashicups" {
  username = "education"
  password = "test123"
  host     = "http://localhost:19090"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `api_url` (String, Sensitive) Url for Cleura API. May also be provided via CLEURA_URL environment variable.
- `domain_id` (String, Sensitive) DomainId for Cleura API. May also be provided via CLEURA_DOMAIN_ID environment variable.
- `password` (String) Password for Cleura API. May also be provided via CLEURA_PW environment variable.
- `username` (String) Username for Cleura API. May also be provided via CLEURA_USER environment variable.