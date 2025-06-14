---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "formancecloud_organization_member Resource - formancecloud"
subcategory: ""
description: |-
  Manages organization members and invitations in Formance Cloud. This resource can be used to invite users to an organization and manage their access levels.
---

# formancecloud_organization_member (Resource)

Manages organization members and invitations in Formance Cloud. This resource can be used to invite users to an organization and manage their access levels.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `email` (String) The email address of the user to invite or add to the organization.
- `organization_id` (String) The organization ID where the member will be added.

### Optional

- `role` (String) The role to assign to the user in the organization. Valid values are: NONE, READ, WRITE.

### Read-Only

- `id` (String) The unique identifier of the invitation or membership.
- `user_id` (String) The user ID once the invitation has been accepted.
