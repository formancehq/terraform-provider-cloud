---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "formancecloud_region_versions Data Source - formancecloud"
subcategory: ""
description: |-
  Retrieves the list of available Formance versions for a specific region.
---

# formancecloud_region_versions (Data Source)

Retrieves the list of available Formance versions for a specific region.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The unique identifier of the region.
- `organization_id` (String) The organization ID that owns the region.

### Read-Only

- `versions` (Attributes List) The list of available Formance versions in the region. (see [below for nested schema](#nestedatt--versions))

<a id="nestedatt--versions"></a>
### Nested Schema for `versions`

Read-Only:

- `name` (String) The version name (e.g., v1.0.0, v2.0.0).
