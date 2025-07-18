---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cloud_stacks Data Source - cloud"
subcategory: ""
description: |-
  Retrieves information about a Formance Cloud stack. If id is specified, returns a specific stack by ID. Otherwise, returns the first available stack sorted alphabetically by name for predictable behavior.
---

# cloud_stacks (Data Source)

Retrieves information about a Formance Cloud stack. If id is specified, returns a specific stack by ID. Otherwise, returns the first available stack sorted alphabetically by name for predictable behavior.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (String) The unique identifier of the stack. If not specified, returns the first available stack sorted alphabetically by name.
- `name` (String) The name of the stack.

### Read-Only

- `region_id` (String) The region ID where the stack is installed.
- `state` (String) The current state of the stack.
- `status` (String) The current status of the stack.
