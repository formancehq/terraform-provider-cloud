# Formance Cloud Terraform Provider

Le provider Terraform Formance Cloud vous permet de gérer vos ressources Formance Cloud via Infrastructure as Code (IaC). Ce provider prend en charge la gestion des organisations, des stacks, des régions et des modules.

## Table des matières

- [Installation](#installation)
- [Configuration](#configuration)
- [Guide de démarrage rapide](#guide-de-démarrage-rapide)
- [Authentification](#authentification)
- [Ressources disponibles](#ressources-disponibles)
- [Data Sources](#data-sources)
- [Exemples](#exemples)
- [Documentation complète](#documentation-complète)
- [Support](#support)

## Installation

### Terraform 0.13+

```hcl
terraform {
  required_providers {
    formancecloud = {
      source  = "formancehq/formancecloud"
      version = "~> 1.0"
    }
  }
}

provider "formancecloud" {
  # Configuration...
}
```

## Configuration

Le provider peut être configuré de deux manières :

### 1. Configuration directe

```hcl
provider "formancecloud" {
  client_id     = "votre-client-id"
  client_secret = "votre-client-secret"
  endpoint      = "https://api.cloud.formance.com" # Optionnel
}
```

### 2. Variables d'environnement

```bash
export FORMANCE_CLOUD_CLIENT_ID="votre-client-id"
export FORMANCE_CLOUD_CLIENT_SECRET="votre-client-secret"
export FORMANCE_CLOUD_API_ENDPOINT="https://api.cloud.formance.com" # Optionnel
```

## Guide de démarrage rapide

Voici un exemple minimal pour démarrer avec le provider Formance Cloud :

```hcl
# Configuration du provider
provider "formancecloud" {
  # Les credentials peuvent être définis via variables d'environnement
}

# Créer une organisation
resource "formancecloud_organization" "main" {
  name = "mon-organisation"
}

# Créer une région privée
resource "formancecloud_region" "europe" {
  name            = "europe-west"
  organization_id = formancecloud_organization.main.id
}

# Créer un stack
resource "formancecloud_stack" "production" {
  name            = "production"
  organization_id = formancecloud_organization.main.id
  region_id       = formancecloud_region.europe.id
}

# Activer le module ledger
resource "formancecloud_stack_module" "ledger" {
  name            = "ledger"
  stack_id        = formancecloud_stack.production.id
  organization_id = formancecloud_organization.main.id
}
```

## Authentification

### Obtenir vos credentials

Le provider utilise l'authentification OAuth2 avec des client credentials. Pour obtenir vos credentials :

1. Connectez-vous à votre compte Formance Cloud
2. Accédez aux paramètres de votre organisation
3. Créez une nouvelle application OAuth2
4. Notez le `client_id` et le `client_secret`

### Bonnes pratiques de sécurité

- **Ne jamais commiter vos credentials** dans votre code
- Utilisez des variables d'environnement ou un gestionnaire de secrets
- Limitez les permissions de vos credentials au strict nécessaire
- Faites tourner régulièrement vos secrets

## Ressources disponibles

### Organizations
- `formancecloud_organization` - Gère une organisation Formance Cloud

### Stacks
- `formancecloud_stack` - Gère un environnement isolé pour vos services Formance

### Regions
- `formancecloud_region` - Gère une région privée dédiée

### Modules
- `formancecloud_stack_module` - Active/désactive des modules sur un stack

### Gestion des accès
- `formancecloud_organization_member` - Gère les membres d'une organisation
- `formancecloud_stack_member` - Gère les accès aux stacks

## Data Sources

- `formancecloud_organizations` - Récupère les informations d'une organisation
- `formancecloud_stacks` - Récupère les informations d'un stack
- `formancecloud_regions` - Récupère les informations d'une région
- `formancecloud_region_versions` - Liste les versions disponibles dans une région

## Exemples

### Déploiement multi-environnements

```hcl
# Variables pour les environnements
variable "environments" {
  default = ["development", "staging", "production"]
}

# Créer un stack pour chaque environnement
resource "formancecloud_stack" "env" {
  for_each        = toset(var.environments)
  name            = each.value
  organization_id = formancecloud_organization.main.id
  region_id       = formancecloud_region.europe.id
}

# Activer les modules nécessaires pour chaque stack
resource "formancecloud_stack_module" "ledger" {
  for_each        = formancecloud_stack.env
  name            = "ledger"
  stack_id        = each.value.id
  organization_id = formancecloud_organization.main.id
}
```

### Gestion des accès avec équipes

```hcl
# Définir les équipes et leurs accès
locals {
  teams = {
    developers = {
      members = ["dev1@example.com", "dev2@example.com"]
      role    = "WRITE"
    }
    observers = {
      members = ["observer1@example.com", "observer2@example.com"]
      role    = "READ"
    }
  }
}

# Ajouter les membres à l'organisation
resource "formancecloud_organization_member" "members" {
  for_each        = toset(flatten([for team in local.teams : team.members]))
  organization_id = formancecloud_organization.main.id
  email          = each.value
  role           = "READ" # Accès minimum à l'organisation
}

# Accorder les accès aux stacks selon les équipes
resource "formancecloud_stack_member" "team_access" {
  for_each = {
    for member in flatten([
      for team_name, team in local.teams : [
        for email in team.members : {
          key     = "${team_name}-${email}"
          email   = email
          role    = team.role
          user_id = formancecloud_organization_member.members[email].user_id
        }
      ]
    ]) : member.key => member
  }
  
  organization_id = formancecloud_organization.main.id
  stack_id       = formancecloud_stack.production.id
  user_id        = each.value.user_id
  role           = each.value.role
}
```

## Documentation complète

Pour plus d'informations détaillées sur chaque ressource et data source :

- [Documentation des ressources](./docs/resources/)
- [Documentation des data sources](./docs/data-sources/)
- [Exemples complets](./examples/)

## Modules disponibles

Les modules suivants peuvent être activés sur vos stacks :

- **ledger** - Moteur comptable central
- **payments** - Gestion et orchestration des paiements
- **webhooks** - Gestion et distribution de webhooks
- **wallets** - Fonctionnalités de portefeuilles numériques
- **search** - Capacités de recherche plein texte
- **reconciliation** - Réconciliation de transactions
- **orchestration** - Orchestration de workflows
- **auth** - Authentification et autorisation
- **stargate** - Gateway API

## Dépannage

### Erreurs courantes

#### Erreur d'authentification
```
Error: Failed to authenticate with Formance Cloud API
```
**Solution** : Vérifiez vos `client_id` et `client_secret`. Assurez-vous qu'ils sont correctement configurés.

#### Erreur de permissions
```
Error: Insufficient permissions to perform this action
```
**Solution** : Vérifiez que vos credentials ont les permissions nécessaires pour l'action demandée.

#### Stack non supprimable
```
Error: Stack cannot be deleted as it contains data
```
**Solution** : Utilisez `force_destroy = true` avec précaution pour forcer la suppression.

## Support

- **Issues GitHub** : [github.com/formancehq/terraform-provider-cloud/issues](https://github.com/formancehq/terraform-provider-cloud/issues)
- **Documentation API** : [docs.formance.com](https://docs.formance.com)
- **Contact** : support@formance.com

## Contribution

Les contributions sont les bienvenues ! Consultez notre [guide de contribution](CONTRIBUTING.md) pour plus d'informations.

## Licence

Ce provider est distribué sous licence Apache 2.0. Voir [LICENSE](LICENSE) pour plus de détails.
