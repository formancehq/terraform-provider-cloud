---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cloud_stack_module Resource - cloud"
subcategory: ""
description: |-
  Manages modules within a Formance Cloud stack. Modules are individual services that can be enabled or disabled on a stack.
---

# cloud_stack_module (Resource)

Manages modules within a Formance Cloud stack. Modules are individual services that can be enabled or disabled on a stack.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the module to enable. Valid module names include: ledger, payments, webhooks, wallets, search, reconciliation, orchestration, auth, stargate.
- `stack_id` (String) The ID of the stack where the module will be enabled.
