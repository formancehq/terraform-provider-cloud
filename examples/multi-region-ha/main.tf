# Configuration multi-régions haute disponibilité
# Cet exemple montre comment déployer Formance Cloud sur plusieurs régions

terraform {
  required_providers {
    cloud = {
      source  = "formancehq/cloud"
      version = "~> 1.0"
    }
  }
}

provider "cloud" {}

# Variables de configuration
variable "organization_name" {
  description = "Nom de l'organisation"
  type        = string
}

variable "regions" {
  description = "Configuration des régions"
  type = map(object({
    name    = string
    primary = bool
    stacks  = list(string)
  }))
  default = {
    europe = {
      name    = "eu-west-1"
      primary = true
      stacks  = ["production-eu", "staging-eu"]
    }
    us = {
      name    = "us-east-1"
      primary = false
      stacks  = ["production-us", "staging-us"]
    }
    asia = {
      name    = "ap-southeast-1"
      primary = false
      stacks  = ["production-asia"]
    }
  }
}

variable "critical_modules" {
  description = "Modules critiques à déployer sur tous les stacks"
  type        = list(string)
  default     = ["ledger", "payments", "auth", "stargate"]
}

variable "disaster_recovery_config" {
  description = "Configuration pour la reprise d'activité"
  type = object({
    enable_cross_region_backup = bool
    backup_retention_days      = number
  })
  default = {
    enable_cross_region_backup = true
    backup_retention_days      = 30
  }
}

# Organisation principale
resource "cloud_organization" "global" {
  name                        = var.organization_name
  domain                      = "${var.organization_name}.global"
  default_organization_access = "READ"
  default_stack_access        = "NONE"
}

# Création des régions
resource "cloud_region" "regions" {
  for_each = var.regions

  name = each.value.name
}

# Récupération des versions par région
data "cloud_region_versions" "versions" {
  for_each = cloud_region.regions

  id = each.value.id
}

# Local pour déterminer la version à utiliser
locals {
  # Utiliser la même version sur toutes les régions pour la cohérence
  global_version = data.cloud_region_versions.versions["europe"].versions[0].name

  # Création d'une structure plate pour les stacks
  stacks_config = flatten([
    for region_key, region in var.regions : [
      for stack_name in region.stacks : {
        key        = "${region_key}-${stack_name}"
        region_key = region_key
        region_id  = cloud_region.regions[region_key].id
        stack_name = stack_name
        is_primary = region.primary
      }
    ]
  ])
}

# Création des stacks dans chaque région
resource "cloud_stack" "multi_region" {
  for_each = {
    for stack in local.stacks_config : stack.key => stack
  }

  name            = each.value.stack_name
  organization_id = cloud_organization.global.id
  region_id       = each.value.region_id
  version         = local.global_version

  # Les stacks de production ne peuvent pas être supprimés facilement
  force_destroy = !contains(["production"], split("-", each.value.stack_name)[0])

  lifecycle {
    # Créer les nouvelles ressources avant de détruire les anciennes
    create_before_destroy = true
  }
}

# Activation des modules critiques sur tous les stacks
resource "cloud_stack_module" "critical" {
  for_each = {
    for pair in setproduct(keys(cloud_stack.multi_region), var.critical_modules) :
    "${pair[0]}-${pair[1]}" => {
      stack_key = pair[0]
      module    = pair[1]
    }
  }

  name            = each.value.module
  stack_id        = cloud_stack.multi_region[each.value.stack_key].id
  organization_id = cloud_organization.global.id
}

# Modules additionnels pour certains stacks
locals {
  additional_modules = {
    "production-eu"   = ["reconciliation", "webhooks", "search"]
    "production-us"   = ["reconciliation", "webhooks"]
    "production-asia" = ["webhooks"]
  }
}

resource "cloud_stack_module" "additional" {
  for_each = {
    for item in flatten([
      for stack_name, modules in local.additional_modules : [
        for module in modules : {
          key        = "${stack_name}-${module}"
          stack_name = stack_name
          module     = module
        }
      ]
    ]) : item.key => item
  }

  name     = each.value.module
  stack_id = [for k, v in cloud_stack.multi_region : v.id if v.name == each.value.stack_name][0]
}

# Configuration des équipes par région
variable "regional_teams" {
  description = "Équipes régionales"
  type = map(list(object({
    email = string
    role  = string
  })))
  default = {
    europe = [
      { email = "eu-lead@company.com", role = "WRITE" },
      { email = "eu-dev@company.com", role = "WRITE" }
    ]
    us = [
      { email = "us-lead@company.com", role = "WRITE" },
      { email = "us-dev@company.com", role = "WRITE" }
    ]
    asia = [
      { email = "asia-lead@company.com", role = "WRITE" }
    ]
  }
}

# Ajout des membres régionaux
resource "cloud_organization_member" "regional_teams" {
  for_each = {
    for item in flatten([
      for region, members in var.regional_teams : [
        for member in members : {
          key   = "${region}-${member.email}"
          email = member.email
          role  = member.role
        }
      ]
    ]) : item.key => item
  }

  organization_id = cloud_organization.global.id
  email           = each.value.email
  role            = each.value.role
}

# Attribution des accès par région
resource "cloud_stack_member" "regional_access" {
  for_each = {
    for item in flatten([
      for region, members in var.regional_teams : [
        for member in members : [
          for stack in cloud_stack.multi_region : {
            key        = "${region}-${member.email}-${stack.name}"
            stack_id   = stack.id
            user_email = member.email
            role       = member.role
          } if contains(var.regions[region].stacks, stack.name)
        ]
      ]
    ]) : item.key => item
  }

  organization_id = cloud_organization.global.id
  stack_id        = each.value.stack_id
  user_id         = cloud_organization_member.regional_teams["${split("-", each.value.key)[0]}-${each.value.user_email}"].user_id
  role            = each.value.role
}

# Stack de disaster recovery
resource "cloud_stack" "dr" {
  name          = "disaster-recovery"
  region_id     = cloud_region.regions["us"].id # DR dans une région différente du primaire
  version       = local.global_version
  force_destroy = false
}

# Modules minimaux pour le DR
resource "cloud_stack_module" "dr_modules" {
  for_each = toset(["ledger", "auth", "stargate"])

  name     = each.value
  stack_id = cloud_stack.dr.id
}

# Outputs pour le monitoring et la configuration
output "regional_endpoints" {
  description = "Endpoints par région"
  value = {
    for key, region in cloud_region.regions : key => {
      region_id = region.id
      base_url  = region.base_url
      stacks = {
        for stack_key, stack in cloud_stack.multi_region :
        stack.name => stack.uri if contains(split("-", stack_key), key)
      }
    }
  }
}

output "primary_region" {
  description = "Région primaire"
  value       = [for k, v in var.regions : k if v.primary][0]
}

output "dr_stack_uri" {
  description = "URI du stack de disaster recovery"
  value       = cloud_stack.dr.uri
  sensitive   = true
}

output "health_check_endpoints" {
  description = "Endpoints pour les health checks"
  value = {
    for key, stack in cloud_stack.multi_region :
    stack.name => "${stack.uri}/health"
  }
}

# Note sur la haute disponibilité
output "ha_configuration_note" {
  value = <<-EOT
    Configuration haute disponibilité déployée:
    - ${length(var.regions)} régions configurées
    - ${length(cloud_stack.multi_region)} stacks déployés
    - Région primaire: ${[for k, v in var.regions : k if v.primary][0]}
    - Stack DR: ${cloud_stack.dr.name}
    
    Pour activer le basculement automatique, configurez votre solution de load balancing
    pour router le trafic entre les régions selon vos besoins.
  EOT
}