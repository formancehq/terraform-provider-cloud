# Exemple complet de configuration Formance Cloud avec Terraform
# Cet exemple montre comment configurer une infrastructure complète

terraform {
  required_providers {
    formancecloud = {
      source  = "formancehq/formancecloud"
      version = "~> 1.0"
    }
  }
}

# Configuration du provider (credentials via variables d'environnement)
provider "formancecloud" {}

# Variables de configuration
variable "organization_name" {
  description = "Nom de l'organisation"
  type        = string
  default     = "my-company"
}

variable "domain" {
  description = "Domaine de l'organisation"
  type        = string
  default     = "mycompany.com"
}

variable "environments" {
  description = "Liste des environnements à créer"
  type        = list(string)
  default     = ["development", "staging", "production"]
}

variable "team_members" {
  description = "Membres de l'équipe et leurs rôles"
  type = map(object({
    email = string
    role  = string
  }))
  default = {
    lead_dev = {
      email = "lead@mycompany.com"
      role  = "WRITE"
    }
    dev1 = {
      email = "dev1@mycompany.com"
      role  = "WRITE"
    }
    dev2 = {
      email = "dev2@mycompany.com"
      role  = "WRITE"
    }
    analyst = {
      email = "analyst@mycompany.com"
      role  = "READ"
    }
  }
}

variable "modules_to_enable" {
  description = "Modules à activer sur les stacks"
  type        = list(string)
  default = [
    "ledger",
    "payments",
    "webhooks",
    "wallets",
    "auth",
    "stargate"
  ]
}

# Création de l'organisation
resource "formancecloud_organization" "main" {
  name                        = var.organization_name
  domain                      = var.domain
  default_organization_access = "READ"
  default_stack_access        = "NONE" # Accès explicite requis
}

# Création d'une région privée pour l'Europe
resource "formancecloud_region" "europe" {
  name            = "europe-west"
  organization_id = formancecloud_organization.main.id
}

# Création d'une région privée pour les US (optionnel)
resource "formancecloud_region" "us" {
  name            = "us-east"
  organization_id = formancecloud_organization.main.id
}

# Récupération des versions disponibles
data "formancecloud_region_versions" "europe" {
  id              = formancecloud_region.europe.id
  organization_id = formancecloud_organization.main.id
}

# Création des stacks pour chaque environnement
resource "formancecloud_stack" "environments" {
  for_each = toset(var.environments)

  name            = each.value
  organization_id = formancecloud_organization.main.id
  region_id       = formancecloud_region.europe.id

  # Utiliser la dernière version stable pour dev/staging, version fixe pour prod
  version = each.value == "production" ? "v2.0.0" : data.formancecloud_region_versions.europe.versions[0].name

  # Protection contre la suppression accidentelle en production
  force_destroy = each.value != "production"

  lifecycle {
    # Empêcher la suppression accidentelle du stack de production
    prevent_destroy = false # Mettre à true en production réelle
  }
}

# Activation des modules sur chaque stack
resource "formancecloud_stack_module" "modules" {
  for_each = {
    for pair in setproduct(keys(formancecloud_stack.environments), var.modules_to_enable) :
    "${pair[0]}-${pair[1]}" => {
      stack_key = pair[0]
      module    = pair[1]
    }
  }

  name            = each.value.module
  stack_id        = formancecloud_stack.environments[each.value.stack_key].id
  organization_id = formancecloud_organization.main.id

  # Les modules ont des dépendances, s'assurer qu'ils sont créés dans le bon ordre
  depends_on = [
    formancecloud_stack.environments
  ]
}

# Ajout des membres à l'organisation
resource "formancecloud_organization_member" "team" {
  for_each = var.team_members

  organization_id = formancecloud_organization.main.id
  email           = each.value.email
  role            = each.value.role
}

# Configuration des accès aux stacks
locals {
  # Matrice des accès : qui a accès à quel environnement
  stack_access = {
    # Tout le monde a accès au dev
    development = {
      for name, member in var.team_members : name => member.role
    }
    # Seuls les devs ont accès au staging
    staging = {
      for name, member in var.team_members : name => member.role
      if member.role == "WRITE"
    }
    # Accès restreint à la production
    production = {
      lead_dev = "WRITE"
      analyst  = "READ"
    }
  }
}

# Attribution des accès aux stacks
resource "formancecloud_stack_member" "access" {
  for_each = {
    for item in flatten([
      for env, access in local.stack_access : [
        for member_name, role in access : {
          key         = "${env}-${member_name}"
          env         = env
          member_name = member_name
          role        = role
        }
      ]
    ]) : item.key => item
  }

  organization_id = formancecloud_organization.main.id
  stack_id        = formancecloud_stack.environments[each.value.env].id
  user_id         = formancecloud_organization_member.team[each.value.member_name].user_id
  role            = each.value.role
}

# Stack dédié pour les tests d'intégration (CI/CD)
resource "formancecloud_stack" "ci" {
  name            = "ci-testing"
  organization_id = formancecloud_organization.main.id
  region_id       = formancecloud_region.europe.id
  force_destroy   = true # Peut être supprimé sans confirmation
}

# Modules minimaux pour les tests CI
resource "formancecloud_stack_module" "ci_modules" {
  for_each = toset(["ledger", "auth"])

  name            = each.value
  stack_id        = formancecloud_stack.ci.id
  organization_id = formancecloud_organization.main.id
}

# Outputs utiles
output "organization_id" {
  description = "ID de l'organisation créée"
  value       = formancecloud_organization.main.id
}

output "stack_urls" {
  description = "URLs des stacks créés"
  value = {
    for name, stack in formancecloud_stack.environments : name => stack.uri
  }
}

output "region_endpoints" {
  description = "Endpoints des régions"
  value = {
    europe = formancecloud_region.europe.base_url
    us     = formancecloud_region.us.base_url
  }
}

output "ci_stack_url" {
  description = "URL du stack CI pour les tests automatisés"
  value       = formancecloud_stack.ci.uri
}

# Note importante sur le secret de région (affiché uniquement à la création)
output "region_secret_note" {
  description = "Note sur les secrets de région"
  value       = "Les secrets de région sont disponibles uniquement lors de la création. Stockez-les de manière sécurisée."
}