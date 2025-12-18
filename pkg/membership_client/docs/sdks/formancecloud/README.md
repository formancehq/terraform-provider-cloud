# FormanceCloud SDK

## Overview

### Available Operations

* [GetServerInfo](#getserverinfo) - Get server info
* [ListOrganizations](#listorganizations) - List organizations of the connected user
* [CreateOrganization](#createorganization) - Create organization
* [~~ListOrganizationsExpanded~~](#listorganizationsexpanded) - List organizations of the connected user with expanded data :warning: **Deprecated**
* [ReadOrganization](#readorganization) - Read organization
* [UpdateOrganization](#updateorganization) - Update organization
* [DeleteOrganization](#deleteorganization) - Delete organization
* [ReadAuthenticationProvider](#readauthenticationprovider) - Read authentication provider
* [UpsertAuthenticationProvider](#upsertauthenticationprovider) - Upsert an authentication provider
* [DeleteAuthenticationProvider](#deleteauthenticationprovider) - Delete authentication provider
* [ListFeatures](#listfeatures) - List features
* [AddFeatures](#addfeatures) - Add Features
* [DeleteFeature](#deletefeature) - Delete feature
* [~~ReadOrganizationClient~~](#readorganizationclient) - Read organization client (DEPRECATED) (until 12/31/2025) :warning: **Deprecated**
* [~~CreateOrganizationClient~~](#createorganizationclient) - Create organization client (DEPRECATED) (until 12/31/2025) :warning: **Deprecated**
* [~~DeleteOrganizationClient~~](#deleteorganizationclient) - Delete organization client (DEPRECATED) (until 12/31/2025) :warning: **Deprecated**
* [OrganizationClientsRead](#organizationclientsread) - Read organization clients
* [OrganizationClientCreate](#organizationclientcreate) - Create organization client
* [OrganizationClientRead](#organizationclientread) - Read organization client
* [OrganizationClientDelete](#organizationclientdelete) - Delete organization client
* [OrganizationClientUpdate](#organizationclientupdate) - Update organization client
* [ListLogs](#listlogs) - List logs
* [ListUsersOfOrganization](#listusersoforganization) - List users of organization
* [ReadUserOfOrganization](#readuseroforganization) - Read user of organization
* [UpsertOrganizationUser](#upsertorganizationuser) - Update user within an organization
* [DeleteUserFromOrganization](#deleteuserfromorganization) - delete user from organization
* [ListPolicies](#listpolicies) - List policies of organization
* [CreatePolicy](#createpolicy) - Create policy
* [ReadPolicy](#readpolicy) - Read policy with scopes
* [UpdatePolicy](#updatepolicy) - Update policy
* [DeletePolicy](#deletepolicy) - Delete policy
* [AddScopeToPolicy](#addscopetopolicy) - Add scope to policy
* [RemoveScopeFromPolicy](#removescopefrompolicy) - Remove scope from policy
* [ListStacks](#liststacks) - List stacks
* [CreateStack](#createstack) - Create stack
* [ListModules](#listmodules) - List modules of a stack
* [EnableModule](#enablemodule) - enable module
* [DisableModule](#disablemodule) - disable module
* [UpgradeStack](#upgradestack) - Upgrade stack
* [GetStack](#getstack) - Find stack
* [UpdateStack](#updatestack) - Update stack
* [DeleteStack](#deletestack) - Delete stack
* [ListStackUsersAccesses](#liststackusersaccesses) - List stack users accesses within an organization
* [ReadStackUserAccess](#readstackuseraccess) - Read stack user access within an organization
* [DeleteStackUserAccess](#deletestackuseraccess) - Delete stack user access within an organization
* [UpsertStackUserAccess](#upsertstackuseraccess) - Update stack user access within an organization
* [DisableStack](#disablestack) - Disable stack
* [EnableStack](#enablestack) - Enable stack
* [RestoreStack](#restorestack) - Restore stack
* [EnableStargate](#enablestargate) - Enable stargate on a stack
* [DisableStargate](#disablestargate) - Disable stargate on a stack
* [ListInvitations](#listinvitations) - List invitations of the user
* [AcceptInvitation](#acceptinvitation) - Accept invitation
* [DeclineInvitation](#declineinvitation) - Decline invitation
* [ListOrganizationInvitations](#listorganizationinvitations) - List invitations of the organization
* [CreateInvitation](#createinvitation) - Create invitation
* [DeleteInvitation](#deleteinvitation) - Delete invitation
* [ListRegions](#listregions) - List regions
* [CreatePrivateRegion](#createprivateregion) - Create a private region
* [GetRegion](#getregion) - Get region
* [DeleteRegion](#deleteregion) - Delete region
* [GetRegionVersions](#getregionversions) - Get region versions
* [ListOrganizationApplications](#listorganizationapplications) - List applications enabled for organization
* [GetOrganizationApplication](#getorganizationapplication) - Get application for organization
* [EnableApplicationForOrganization](#enableapplicationfororganization) - Enable application for organization
* [DisableApplicationForOrganization](#disableapplicationfororganization) - Disable application for organization
* [ListApplications](#listapplications) - List applications
* [CreateApplication](#createapplication) - Create application
* [GetApplication](#getapplication) - Get application
* [UpdateApplication](#updateapplication) - Update application
* [DeleteApplication](#deleteapplication) - Delete application
* [CreateApplicationScope](#createapplicationscope) - Create application scope
* [DeleteApplicationScope](#deleteapplicationscope) - Delete application scope
* [CreateUser](#createuser) - Create user
* [ReadConnectedUser](#readconnecteduser) - Read user

## GetServerInfo

Get server info

### Example Usage

<!-- UsageSnippet language="go" operationID="getServerInfo" method="get" path="/_info" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.GetServerInfo(ctx)
    if err != nil {
        log.Fatal(err)
    }
    if res.ServerInfo != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.GetServerInfoResponse](../../pkg/models/operations/getserverinforesponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ListOrganizations

List organizations of the connected user

### Example Usage

<!-- UsageSnippet language="go" operationID="listOrganizations" method="get" path="/organizations" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ListOrganizations(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.ListOrganizationExpandedResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `expand`                                                     | **bool*                                                      | :heavy_minus_sign:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ListOrganizationsResponse](../../pkg/models/operations/listorganizationsresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## CreateOrganization

Create organization

### Example Usage

<!-- UsageSnippet language="go" operationID="createOrganization" method="post" path="/organizations" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.CreateOrganization(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.CreateOrganizationResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                | Type                                                                                     | Required                                                                                 | Description                                                                              |
| ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `ctx`                                                                                    | [context.Context](https://pkg.go.dev/context#Context)                                    | :heavy_check_mark:                                                                       | The context to use for the request.                                                      |
| `request`                                                                                | [shared.CreateOrganizationRequest](../../pkg/models/shared/createorganizationrequest.md) | :heavy_check_mark:                                                                       | The request object to use for the request.                                               |
| `opts`                                                                                   | [][operations.Option](../../pkg/models/operations/option.md)                             | :heavy_minus_sign:                                                                       | The options for this request.                                                            |

### Response

**[*operations.CreateOrganizationResponse](../../pkg/models/operations/createorganizationresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ~~ListOrganizationsExpanded~~

List organizations of the connected user with expanded data

> :warning: **DEPRECATED**: This will be removed in a future release, please migrate away from it as soon as possible.

### Example Usage

<!-- UsageSnippet language="go" operationID="listOrganizationsExpanded" method="get" path="/organizations/expanded" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ListOrganizationsExpanded(ctx)
    if err != nil {
        log.Fatal(err)
    }
    if res.ListOrganizationExpandedResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ListOrganizationsExpandedResponse](../../pkg/models/operations/listorganizationsexpandedresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ReadOrganization

Read organization

### Example Usage

<!-- UsageSnippet language="go" operationID="readOrganization" method="get" path="/organizations/{organizationId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ReadOrganization(ctx, "<id>", nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.ReadOrganizationResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `expand`                                                     | **bool*                                                      | :heavy_minus_sign:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ReadOrganizationResponse](../../pkg/models/operations/readorganizationresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## UpdateOrganization

Update organization

### Example Usage

<!-- UsageSnippet language="go" operationID="updateOrganization" method="put" path="/organizations/{organizationId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.UpdateOrganization(ctx, "<id>", nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.ReadOrganizationResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                               | Type                                                                    | Required                                                                | Description                                                             |
| ----------------------------------------------------------------------- | ----------------------------------------------------------------------- | ----------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| `ctx`                                                                   | [context.Context](https://pkg.go.dev/context#Context)                   | :heavy_check_mark:                                                      | The context to use for the request.                                     |
| `organizationID`                                                        | *string*                                                                | :heavy_check_mark:                                                      | N/A                                                                     |
| `organizationData`                                                      | [*shared.OrganizationData](../../pkg/models/shared/organizationdata.md) | :heavy_minus_sign:                                                      | N/A                                                                     |
| `opts`                                                                  | [][operations.Option](../../pkg/models/operations/option.md)            | :heavy_minus_sign:                                                      | The options for this request.                                           |

### Response

**[*operations.UpdateOrganizationResponse](../../pkg/models/operations/updateorganizationresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## DeleteOrganization

Delete organization

### Example Usage

<!-- UsageSnippet language="go" operationID="deleteOrganization" method="delete" path="/organizations/{organizationId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.DeleteOrganization(ctx, "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.DeleteOrganizationResponse](../../pkg/models/operations/deleteorganizationresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ReadAuthenticationProvider

Read authentication provider

### Example Usage

<!-- UsageSnippet language="go" operationID="readAuthenticationProvider" method="get" path="/organizations/{organizationId}/authentication-provider" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ReadAuthenticationProvider(ctx, "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.AuthenticationProviderResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ReadAuthenticationProviderResponse](../../pkg/models/operations/readauthenticationproviderresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## UpsertAuthenticationProvider

Upsert an authentication provider

### Example Usage

<!-- UsageSnippet language="go" operationID="upsertAuthenticationProvider" method="put" path="/organizations/{organizationId}/authentication-provider" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.UpsertAuthenticationProvider(ctx, "<id>", nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.AuthenticationProviderResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                   | Type                                                                                        | Required                                                                                    | Description                                                                                 |
| ------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------- |
| `ctx`                                                                                       | [context.Context](https://pkg.go.dev/context#Context)                                       | :heavy_check_mark:                                                                          | The context to use for the request.                                                         |
| `organizationID`                                                                            | *string*                                                                                    | :heavy_check_mark:                                                                          | N/A                                                                                         |
| `authenticationProviderData`                                                                | [*shared.AuthenticationProviderData](../../pkg/models/shared/authenticationproviderdata.md) | :heavy_minus_sign:                                                                          | N/A                                                                                         |
| `opts`                                                                                      | [][operations.Option](../../pkg/models/operations/option.md)                                | :heavy_minus_sign:                                                                          | The options for this request.                                                               |

### Response

**[*operations.UpsertAuthenticationProviderResponse](../../pkg/models/operations/upsertauthenticationproviderresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## DeleteAuthenticationProvider

Delete authentication provider

### Example Usage

<!-- UsageSnippet language="go" operationID="deleteAuthenticationProvider" method="delete" path="/organizations/{organizationId}/authentication-provider" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.DeleteAuthenticationProvider(ctx, "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.DeleteAuthenticationProviderResponse](../../pkg/models/operations/deleteauthenticationproviderresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ListFeatures

List features

### Example Usage

<!-- UsageSnippet language="go" operationID="listFeatures" method="get" path="/organizations/{organizationId}/features" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ListFeatures(ctx, "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Object != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ListFeaturesResponse](../../pkg/models/operations/listfeaturesresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## AddFeatures

Add Features

### Example Usage

<!-- UsageSnippet language="go" operationID="addFeatures" method="post" path="/organizations/{organizationId}/features" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.AddFeatures(ctx, "<id>", nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                   | Type                                                                                        | Required                                                                                    | Description                                                                                 |
| ------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------- |
| `ctx`                                                                                       | [context.Context](https://pkg.go.dev/context#Context)                                       | :heavy_check_mark:                                                                          | The context to use for the request.                                                         |
| `organizationID`                                                                            | *string*                                                                                    | :heavy_check_mark:                                                                          | N/A                                                                                         |
| `requestBody`                                                                               | [*operations.AddFeaturesRequestBody](../../pkg/models/operations/addfeaturesrequestbody.md) | :heavy_minus_sign:                                                                          | N/A                                                                                         |
| `opts`                                                                                      | [][operations.Option](../../pkg/models/operations/option.md)                                | :heavy_minus_sign:                                                                          | The options for this request.                                                               |

### Response

**[*operations.AddFeaturesResponse](../../pkg/models/operations/addfeaturesresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## DeleteFeature

Delete feature

### Example Usage

<!-- UsageSnippet language="go" operationID="deleteFeature" method="delete" path="/organizations/{organizationId}/features/{name}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.DeleteFeature(ctx, "<id>", "<value>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `name`                                                       | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.DeleteFeatureResponse](../../pkg/models/operations/deletefeatureresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ~~ReadOrganizationClient~~

Read organization client (DEPRECATED) (until 12/31/2025)

> :warning: **DEPRECATED**: This will be removed in a future release, please migrate away from it as soon as possible.

### Example Usage

<!-- UsageSnippet language="go" operationID="readOrganizationClient" method="get" path="/organizations/{organizationId}/client" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ReadOrganizationClient(ctx, "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.CreateClientResponseResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ReadOrganizationClientResponse](../../pkg/models/operations/readorganizationclientresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ~~CreateOrganizationClient~~

Create organization client (DEPRECATED) (until 12/31/2025)

> :warning: **DEPRECATED**: This will be removed in a future release, please migrate away from it as soon as possible.

### Example Usage

<!-- UsageSnippet language="go" operationID="createOrganizationClient" method="put" path="/organizations/{organizationId}/client" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.CreateOrganizationClient(ctx, "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.CreateClientResponseResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.CreateOrganizationClientResponse](../../pkg/models/operations/createorganizationclientresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ~~DeleteOrganizationClient~~

Delete organization client (DEPRECATED) (until 12/31/2025)

> :warning: **DEPRECATED**: This will be removed in a future release, please migrate away from it as soon as possible.

### Example Usage

<!-- UsageSnippet language="go" operationID="deleteOrganizationClient" method="delete" path="/organizations/{organizationId}/client" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.DeleteOrganizationClient(ctx, "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.DeleteOrganizationClientResponse](../../pkg/models/operations/deleteorganizationclientresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## OrganizationClientsRead

Read organization clients

### Example Usage

<!-- UsageSnippet language="go" operationID="organizationClientsRead" method="get" path="/organizations/{organizationId}/clients" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.OrganizationClientsRead(ctx, "<id>", nil, nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.ReadOrganizationClientsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `cursor`                                                     | **string*                                                    | :heavy_minus_sign:                                           | N/A                                                          |
| `pageSize`                                                   | **int64*                                                     | :heavy_minus_sign:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.OrganizationClientsReadResponse](../../pkg/models/operations/organizationclientsreadresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## OrganizationClientCreate

Create organization client

### Example Usage

<!-- UsageSnippet language="go" operationID="organizationClientCreate" method="post" path="/organizations/{organizationId}/clients" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.OrganizationClientCreate(ctx, "<id>", nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.CreateOrganizationClientResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                             | Type                                                                                                  | Required                                                                                              | Description                                                                                           |
| ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                 | [context.Context](https://pkg.go.dev/context#Context)                                                 | :heavy_check_mark:                                                                                    | The context to use for the request.                                                                   |
| `organizationID`                                                                                      | *string*                                                                                              | :heavy_check_mark:                                                                                    | N/A                                                                                                   |
| `createOrganizationClientRequest`                                                                     | [*shared.CreateOrganizationClientRequest](../../pkg/models/shared/createorganizationclientrequest.md) | :heavy_minus_sign:                                                                                    | N/A                                                                                                   |
| `opts`                                                                                                | [][operations.Option](../../pkg/models/operations/option.md)                                          | :heavy_minus_sign:                                                                                    | The options for this request.                                                                         |

### Response

**[*operations.OrganizationClientCreateResponse](../../pkg/models/operations/organizationclientcreateresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## OrganizationClientRead

Read organization client

### Example Usage

<!-- UsageSnippet language="go" operationID="organizationClientRead" method="get" path="/organizations/{organizationId}/clients/{clientId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.OrganizationClientRead(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.ReadOrganizationClientResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `clientID`                                                   | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.OrganizationClientReadResponse](../../pkg/models/operations/organizationclientreadresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## OrganizationClientDelete

Delete organization client

### Example Usage

<!-- UsageSnippet language="go" operationID="organizationClientDelete" method="delete" path="/organizations/{organizationId}/clients/{clientId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.OrganizationClientDelete(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `clientID`                                                   | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.OrganizationClientDeleteResponse](../../pkg/models/operations/organizationclientdeleteresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## OrganizationClientUpdate

Update organization client

### Example Usage

<!-- UsageSnippet language="go" operationID="organizationClientUpdate" method="put" path="/organizations/{organizationId}/clients/{clientId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.OrganizationClientUpdate(ctx, "<id>", "<id>", nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                             | Type                                                                                                  | Required                                                                                              | Description                                                                                           |
| ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                 | [context.Context](https://pkg.go.dev/context#Context)                                                 | :heavy_check_mark:                                                                                    | The context to use for the request.                                                                   |
| `organizationID`                                                                                      | *string*                                                                                              | :heavy_check_mark:                                                                                    | N/A                                                                                                   |
| `clientID`                                                                                            | *string*                                                                                              | :heavy_check_mark:                                                                                    | N/A                                                                                                   |
| `updateOrganizationClientRequest`                                                                     | [*shared.UpdateOrganizationClientRequest](../../pkg/models/shared/updateorganizationclientrequest.md) | :heavy_minus_sign:                                                                                    | N/A                                                                                                   |
| `opts`                                                                                                | [][operations.Option](../../pkg/models/operations/option.md)                                          | :heavy_minus_sign:                                                                                    | The options for this request.                                                                         |

### Response

**[*operations.OrganizationClientUpdateResponse](../../pkg/models/operations/organizationclientupdateresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ListLogs

List logs

### Example Usage

<!-- UsageSnippet language="go" operationID="listLogs" method="get" path="/organizations/{organizationId}/logs" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/operations"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ListLogs(ctx, operations.ListLogsRequest{
        OrganizationID: "<id>",
    })
    if err != nil {
        log.Fatal(err)
    }
    if res.LogCursor != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                    | Type                                                                         | Required                                                                     | Description                                                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `ctx`                                                                        | [context.Context](https://pkg.go.dev/context#Context)                        | :heavy_check_mark:                                                           | The context to use for the request.                                          |
| `request`                                                                    | [operations.ListLogsRequest](../../pkg/models/operations/listlogsrequest.md) | :heavy_check_mark:                                                           | The request object to use for the request.                                   |
| `opts`                                                                       | [][operations.Option](../../pkg/models/operations/option.md)                 | :heavy_minus_sign:                                                           | The options for this request.                                                |

### Response

**[*operations.ListLogsResponse](../../pkg/models/operations/listlogsresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ListUsersOfOrganization

List users of organization

### Example Usage

<!-- UsageSnippet language="go" operationID="listUsersOfOrganization" method="get" path="/organizations/{organizationId}/users" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ListUsersOfOrganization(ctx, "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.ListUsersResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ListUsersOfOrganizationResponse](../../pkg/models/operations/listusersoforganizationresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ReadUserOfOrganization

Read user of organization

### Example Usage

<!-- UsageSnippet language="go" operationID="readUserOfOrganization" method="get" path="/organizations/{organizationId}/users/{userId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ReadUserOfOrganization(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.ReadOrganizationUserResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `userID`                                                     | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ReadUserOfOrganizationResponse](../../pkg/models/operations/readuseroforganizationresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## UpsertOrganizationUser

Update user within an organization

### Example Usage

<!-- UsageSnippet language="go" operationID="upsertOrganizationUser" method="put" path="/organizations/{organizationId}/users/{userId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.UpsertOrganizationUser(ctx, "<id>", "<id>", nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                         | Type                                                                                              | Required                                                                                          | Description                                                                                       |
| ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                             | [context.Context](https://pkg.go.dev/context#Context)                                             | :heavy_check_mark:                                                                                | The context to use for the request.                                                               |
| `organizationID`                                                                                  | *string*                                                                                          | :heavy_check_mark:                                                                                | N/A                                                                                               |
| `userID`                                                                                          | *string*                                                                                          | :heavy_check_mark:                                                                                | N/A                                                                                               |
| `updateOrganizationUserRequest`                                                                   | [*shared.UpdateOrganizationUserRequest](../../pkg/models/shared/updateorganizationuserrequest.md) | :heavy_minus_sign:                                                                                | N/A                                                                                               |
| `opts`                                                                                            | [][operations.Option](../../pkg/models/operations/option.md)                                      | :heavy_minus_sign:                                                                                | The options for this request.                                                                     |

### Response

**[*operations.UpsertOrganizationUserResponse](../../pkg/models/operations/upsertorganizationuserresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## DeleteUserFromOrganization

The owner of the organization can remove anyone while each user can leave any organization where it is not owner.


### Example Usage

<!-- UsageSnippet language="go" operationID="deleteUserFromOrganization" method="delete" path="/organizations/{organizationId}/users/{userId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.DeleteUserFromOrganization(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `userID`                                                     | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.DeleteUserFromOrganizationResponse](../../pkg/models/operations/deleteuserfromorganizationresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ListPolicies

List policies of organization

### Example Usage

<!-- UsageSnippet language="go" operationID="listPolicies" method="get" path="/organizations/{organizationId}/policies" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ListPolicies(ctx, "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.ListPoliciesResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ListPoliciesResponse](../../pkg/models/operations/listpoliciesresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## CreatePolicy

Create policy

### Example Usage

<!-- UsageSnippet language="go" operationID="createPolicy" method="post" path="/organizations/{organizationId}/policies" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.CreatePolicy(ctx, "<id>", nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.CreatePolicyResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `policyData`                                                 | [*shared.PolicyData](../../pkg/models/shared/policydata.md)  | :heavy_minus_sign:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.CreatePolicyResponse](../../pkg/models/operations/createpolicyresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ReadPolicy

Read policy with scopes

### Example Usage

<!-- UsageSnippet language="go" operationID="readPolicy" method="get" path="/organizations/{organizationId}/policies/{policyId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ReadPolicy(ctx, "<id>", 831591)
    if err != nil {
        log.Fatal(err)
    }
    if res.ReadPolicyResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `policyID`                                                   | *int64*                                                      | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ReadPolicyResponse](../../pkg/models/operations/readpolicyresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## UpdatePolicy

Update policy

### Example Usage

<!-- UsageSnippet language="go" operationID="updatePolicy" method="put" path="/organizations/{organizationId}/policies/{policyId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.UpdatePolicy(ctx, "<id>", 127460, nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.UpdatePolicyResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `policyID`                                                   | *int64*                                                      | :heavy_check_mark:                                           | N/A                                                          |
| `policyData`                                                 | [*shared.PolicyData](../../pkg/models/shared/policydata.md)  | :heavy_minus_sign:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.UpdatePolicyResponse](../../pkg/models/operations/updatepolicyresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## DeletePolicy

Delete policy

### Example Usage

<!-- UsageSnippet language="go" operationID="deletePolicy" method="delete" path="/organizations/{organizationId}/policies/{policyId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.DeletePolicy(ctx, "<id>", 114294)
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `policyID`                                                   | *int64*                                                      | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.DeletePolicyResponse](../../pkg/models/operations/deletepolicyresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## AddScopeToPolicy

Add scope to policy

### Example Usage

<!-- UsageSnippet language="go" operationID="addScopeToPolicy" method="put" path="/organizations/{organizationId}/policies/{policyId}/scopes/{scopeId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.AddScopeToPolicy(ctx, "<id>", 328027, 675877)
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `policyID`                                                   | *int64*                                                      | :heavy_check_mark:                                           | N/A                                                          |
| `scopeID`                                                    | *int64*                                                      | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.AddScopeToPolicyResponse](../../pkg/models/operations/addscopetopolicyresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## RemoveScopeFromPolicy

Remove scope from policy

### Example Usage

<!-- UsageSnippet language="go" operationID="removeScopeFromPolicy" method="delete" path="/organizations/{organizationId}/policies/{policyId}/scopes/{scopeId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.RemoveScopeFromPolicy(ctx, "<id>", 995736, 485996)
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `policyID`                                                   | *int64*                                                      | :heavy_check_mark:                                           | N/A                                                          |
| `scopeID`                                                    | *int64*                                                      | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.RemoveScopeFromPolicyResponse](../../pkg/models/operations/removescopefrompolicyresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ListStacks

List stacks

### Example Usage

<!-- UsageSnippet language="go" operationID="listStacks" method="get" path="/organizations/{organizationId}/stacks" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ListStacks(ctx, "<id>", nil, nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.ListStacksResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                                                                       | Type                                                                                                                                            | Required                                                                                                                                        | Description                                                                                                                                     |
| ----------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                                                                           | [context.Context](https://pkg.go.dev/context#Context)                                                                                           | :heavy_check_mark:                                                                                                                              | The context to use for the request.                                                                                                             |
| `organizationID`                                                                                                                                | *string*                                                                                                                                        | :heavy_check_mark:                                                                                                                              | N/A                                                                                                                                             |
| `all`                                                                                                                                           | **bool*                                                                                                                                         | :heavy_minus_sign:                                                                                                                              | Include deleted and disabled stacks                                                                                                             |
| `deleted`                                                                                                                                       | **bool*                                                                                                                                         | :heavy_minus_sign:                                                                                                                              | : warning: ** DEPRECATED **: This will be removed in a future release, please migrate away from it as soon as possible.<br/><br/>Include deleted stacks |
| `opts`                                                                                                                                          | [][operations.Option](../../pkg/models/operations/option.md)                                                                                    | :heavy_minus_sign:                                                                                                                              | The options for this request.                                                                                                                   |

### Response

**[*operations.ListStacksResponse](../../pkg/models/operations/liststacksresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## CreateStack

Create stack

### Example Usage

<!-- UsageSnippet language="go" operationID="createStack" method="post" path="/organizations/{organizationId}/stacks" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.CreateStack(ctx, "<id>", nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.CreateStackResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                   | Type                                                                        | Required                                                                    | Description                                                                 |
| --------------------------------------------------------------------------- | --------------------------------------------------------------------------- | --------------------------------------------------------------------------- | --------------------------------------------------------------------------- |
| `ctx`                                                                       | [context.Context](https://pkg.go.dev/context#Context)                       | :heavy_check_mark:                                                          | The context to use for the request.                                         |
| `organizationID`                                                            | *string*                                                                    | :heavy_check_mark:                                                          | N/A                                                                         |
| `createStackRequest`                                                        | [*shared.CreateStackRequest](../../pkg/models/shared/createstackrequest.md) | :heavy_minus_sign:                                                          | N/A                                                                         |
| `opts`                                                                      | [][operations.Option](../../pkg/models/operations/option.md)                | :heavy_minus_sign:                                                          | The options for this request.                                               |

### Response

**[*operations.CreateStackResponse](../../pkg/models/operations/createstackresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ListModules

List modules of a stack

### Example Usage

<!-- UsageSnippet language="go" operationID="listModules" method="get" path="/organizations/{organizationId}/stacks/{stackId}/modules" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ListModules(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.ListModulesResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `stackID`                                                    | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ListModulesResponse](../../pkg/models/operations/listmodulesresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## EnableModule

enable module

### Example Usage

<!-- UsageSnippet language="go" operationID="enableModule" method="post" path="/organizations/{organizationId}/stacks/{stackId}/modules" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.EnableModule(ctx, "<id>", "<id>", "<value>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `stackID`                                                    | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `name`                                                       | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.EnableModuleResponse](../../pkg/models/operations/enablemoduleresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## DisableModule

disable module

### Example Usage

<!-- UsageSnippet language="go" operationID="disableModule" method="delete" path="/organizations/{organizationId}/stacks/{stackId}/modules" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.DisableModule(ctx, "<id>", "<id>", "<value>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `stackID`                                                    | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `name`                                                       | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.DisableModuleResponse](../../pkg/models/operations/disablemoduleresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## UpgradeStack

Upgrade stack

### Example Usage

<!-- UsageSnippet language="go" operationID="upgradeStack" method="put" path="/organizations/{organizationId}/stacks/{stackId}/upgrade" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.UpgradeStack(ctx, "<id>", "<id>", nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                       | Type                                                            | Required                                                        | Description                                                     |
| --------------------------------------------------------------- | --------------------------------------------------------------- | --------------------------------------------------------------- | --------------------------------------------------------------- |
| `ctx`                                                           | [context.Context](https://pkg.go.dev/context#Context)           | :heavy_check_mark:                                              | The context to use for the request.                             |
| `organizationID`                                                | *string*                                                        | :heavy_check_mark:                                              | N/A                                                             |
| `stackID`                                                       | *string*                                                        | :heavy_check_mark:                                              | N/A                                                             |
| `stackVersion`                                                  | [*shared.StackVersion](../../pkg/models/shared/stackversion.md) | :heavy_minus_sign:                                              | N/A                                                             |
| `opts`                                                          | [][operations.Option](../../pkg/models/operations/option.md)    | :heavy_minus_sign:                                              | The options for this request.                                   |

### Response

**[*operations.UpgradeStackResponse](../../pkg/models/operations/upgradestackresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## GetStack

Find stack

### Example Usage

<!-- UsageSnippet language="go" operationID="getStack" method="get" path="/organizations/{organizationId}/stacks/{stackId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.GetStack(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.CreateStackResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `stackID`                                                    | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.GetStackResponse](../../pkg/models/operations/getstackresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## UpdateStack

Update stack

### Example Usage

<!-- UsageSnippet language="go" operationID="updateStack" method="put" path="/organizations/{organizationId}/stacks/{stackId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.UpdateStack(ctx, "<id>", "<id>", nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.CreateStackResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `stackID`                                                    | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `stackData`                                                  | [*shared.StackData](../../pkg/models/shared/stackdata.md)    | :heavy_minus_sign:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.UpdateStackResponse](../../pkg/models/operations/updatestackresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## DeleteStack

Delete stack

### Example Usage

<!-- UsageSnippet language="go" operationID="deleteStack" method="delete" path="/organizations/{organizationId}/stacks/{stackId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.DeleteStack(ctx, "<id>", "<id>", nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `stackID`                                                    | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `force`                                                      | **bool*                                                      | :heavy_minus_sign:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.DeleteStackResponse](../../pkg/models/operations/deletestackresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ListStackUsersAccesses

List stack users accesses within an organization

### Example Usage

<!-- UsageSnippet language="go" operationID="listStackUsersAccesses" method="get" path="/organizations/{organizationId}/stacks/{stackId}/users" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ListStackUsersAccesses(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.StackUserAccessResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `stackID`                                                    | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ListStackUsersAccessesResponse](../../pkg/models/operations/liststackusersaccessesresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ReadStackUserAccess

Read stack user access within an organization

### Example Usage

<!-- UsageSnippet language="go" operationID="readStackUserAccess" method="get" path="/organizations/{organizationId}/stacks/{stackId}/users/{userId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ReadStackUserAccess(ctx, "<id>", "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.ReadStackUserAccess != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `stackID`                                                    | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `userID`                                                     | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ReadStackUserAccessResponse](../../pkg/models/operations/readstackuseraccessresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## DeleteStackUserAccess

Delete stack user access within an organization

### Example Usage

<!-- UsageSnippet language="go" operationID="deleteStackUserAccess" method="delete" path="/organizations/{organizationId}/stacks/{stackId}/users/{userId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.DeleteStackUserAccess(ctx, "<id>", "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `stackID`                                                    | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `userID`                                                     | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.DeleteStackUserAccessResponse](../../pkg/models/operations/deletestackuseraccessresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## UpsertStackUserAccess

Update stack user access within an organization

### Example Usage

<!-- UsageSnippet language="go" operationID="upsertStackUserAccess" method="put" path="/organizations/{organizationId}/stacks/{stackId}/users/{userId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.UpsertStackUserAccess(ctx, "<id>", "<id>", "<id>", nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                           | Type                                                                                | Required                                                                            | Description                                                                         |
| ----------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------- |
| `ctx`                                                                               | [context.Context](https://pkg.go.dev/context#Context)                               | :heavy_check_mark:                                                                  | The context to use for the request.                                                 |
| `organizationID`                                                                    | *string*                                                                            | :heavy_check_mark:                                                                  | N/A                                                                                 |
| `stackID`                                                                           | *string*                                                                            | :heavy_check_mark:                                                                  | N/A                                                                                 |
| `userID`                                                                            | *string*                                                                            | :heavy_check_mark:                                                                  | N/A                                                                                 |
| `updateStackUserRequest`                                                            | [*shared.UpdateStackUserRequest](../../pkg/models/shared/updatestackuserrequest.md) | :heavy_minus_sign:                                                                  | N/A                                                                                 |
| `opts`                                                                              | [][operations.Option](../../pkg/models/operations/option.md)                        | :heavy_minus_sign:                                                                  | The options for this request.                                                       |

### Response

**[*operations.UpsertStackUserAccessResponse](../../pkg/models/operations/upsertstackuseraccessresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## DisableStack

Disable stack

### Example Usage

<!-- UsageSnippet language="go" operationID="disableStack" method="put" path="/organizations/{organizationId}/stacks/{stackId}/disable" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.DisableStack(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `stackID`                                                    | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.DisableStackResponse](../../pkg/models/operations/disablestackresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## EnableStack

Enable stack

### Example Usage

<!-- UsageSnippet language="go" operationID="enableStack" method="put" path="/organizations/{organizationId}/stacks/{stackId}/enable" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.EnableStack(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `stackID`                                                    | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.EnableStackResponse](../../pkg/models/operations/enablestackresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## RestoreStack

Restore stack

### Example Usage

<!-- UsageSnippet language="go" operationID="restoreStack" method="put" path="/organizations/{organizationId}/stacks/{stackId}/restore" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.RestoreStack(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.CreateStackResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `stackID`                                                    | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.RestoreStackResponse](../../pkg/models/operations/restorestackresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## EnableStargate

Enable stargate on a stack

### Example Usage

<!-- UsageSnippet language="go" operationID="enableStargate" method="put" path="/organizations/{organizationId}/stacks/{stackId}/stargate/enable" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.EnableStargate(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `stackID`                                                    | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.EnableStargateResponse](../../pkg/models/operations/enablestargateresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## DisableStargate

Disable stargate on a stack

### Example Usage

<!-- UsageSnippet language="go" operationID="disableStargate" method="put" path="/organizations/{organizationId}/stacks/{stackId}/stargate/disable" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.DisableStargate(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `stackID`                                                    | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.DisableStargateResponse](../../pkg/models/operations/disablestargateresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ListInvitations

List invitations of the user

### Example Usage

<!-- UsageSnippet language="go" operationID="listInvitations" method="get" path="/me/invitations" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ListInvitations(ctx, nil, nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.ListInvitationsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `status`                                                     | **string*                                                    | :heavy_minus_sign:                                           | Status of organizations                                      |
| `organization`                                               | **string*                                                    | :heavy_minus_sign:                                           | Status of organizations                                      |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ListInvitationsResponse](../../pkg/models/operations/listinvitationsresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## AcceptInvitation

Accept invitation

### Example Usage

<!-- UsageSnippet language="go" operationID="acceptInvitation" method="post" path="/me/invitations/{invitationId}/accept" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.AcceptInvitation(ctx, "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `invitationID`                                               | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.AcceptInvitationResponse](../../pkg/models/operations/acceptinvitationresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## DeclineInvitation

Decline invitation

### Example Usage

<!-- UsageSnippet language="go" operationID="declineInvitation" method="post" path="/me/invitations/{invitationId}/reject" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.DeclineInvitation(ctx, "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `invitationID`                                               | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.DeclineInvitationResponse](../../pkg/models/operations/declineinvitationresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ListOrganizationInvitations

List invitations of the organization

### Example Usage

<!-- UsageSnippet language="go" operationID="listOrganizationInvitations" method="get" path="/organizations/{organizationId}/invitations" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ListOrganizationInvitations(ctx, "<id>", nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.ListInvitationsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `status`                                                     | **string*                                                    | :heavy_minus_sign:                                           | Status of organizations                                      |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ListOrganizationInvitationsResponse](../../pkg/models/operations/listorganizationinvitationsresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## CreateInvitation

Create invitation

### Example Usage

<!-- UsageSnippet language="go" operationID="createInvitation" method="post" path="/organizations/{organizationId}/invitations" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.CreateInvitation(ctx, "<id>", "Manley_Hoeger@hotmail.com")
    if err != nil {
        log.Fatal(err)
    }
    if res.CreateInvitationResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `email`                                                      | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.CreateInvitationResponse](../../pkg/models/operations/createinvitationresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## DeleteInvitation

Delete invitation

### Example Usage

<!-- UsageSnippet language="go" operationID="deleteInvitation" method="delete" path="/organizations/{organizationId}/invitations/{invitationId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.DeleteInvitation(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `invitationID`                                               | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.DeleteInvitationResponse](../../pkg/models/operations/deleteinvitationresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ListRegions

List regions

### Example Usage

<!-- UsageSnippet language="go" operationID="listRegions" method="get" path="/organizations/{organizationId}/regions" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ListRegions(ctx, "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.ListRegionsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ListRegionsResponse](../../pkg/models/operations/listregionsresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## CreatePrivateRegion

Create a private region

### Example Usage

<!-- UsageSnippet language="go" operationID="createPrivateRegion" method="post" path="/organizations/{organizationId}/regions" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.CreatePrivateRegion(ctx, "<id>", nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.CreatedPrivateRegionResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                   | Type                                                                                        | Required                                                                                    | Description                                                                                 |
| ------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------- |
| `ctx`                                                                                       | [context.Context](https://pkg.go.dev/context#Context)                                       | :heavy_check_mark:                                                                          | The context to use for the request.                                                         |
| `organizationID`                                                                            | *string*                                                                                    | :heavy_check_mark:                                                                          | N/A                                                                                         |
| `createPrivateRegionRequest`                                                                | [*shared.CreatePrivateRegionRequest](../../pkg/models/shared/createprivateregionrequest.md) | :heavy_minus_sign:                                                                          | N/A                                                                                         |
| `opts`                                                                                      | [][operations.Option](../../pkg/models/operations/option.md)                                | :heavy_minus_sign:                                                                          | The options for this request.                                                               |

### Response

**[*operations.CreatePrivateRegionResponse](../../pkg/models/operations/createprivateregionresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## GetRegion

Get region

### Example Usage

<!-- UsageSnippet language="go" operationID="getRegion" method="get" path="/organizations/{organizationId}/regions/{regionID}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.GetRegion(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.GetRegionResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `regionID`                                                   | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.GetRegionResponse](../../pkg/models/operations/getregionresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## DeleteRegion

Delete region

### Example Usage

<!-- UsageSnippet language="go" operationID="deleteRegion" method="delete" path="/organizations/{organizationId}/regions/{regionID}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.DeleteRegion(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `regionID`                                                   | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.DeleteRegionResponse](../../pkg/models/operations/deleteregionresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## GetRegionVersions

Get region versions

### Example Usage

<!-- UsageSnippet language="go" operationID="getRegionVersions" method="get" path="/organizations/{organizationId}/regions/{regionID}/versions" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.GetRegionVersions(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.GetRegionVersionsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `regionID`                                                   | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.GetRegionVersionsResponse](../../pkg/models/operations/getregionversionsresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ListOrganizationApplications

List applications enabled for organization

### Example Usage

<!-- UsageSnippet language="go" operationID="listOrganizationApplications" method="get" path="/organizations/{organizationId}/applications" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ListOrganizationApplications(ctx, "<id>", membershipclient.Pointer[int64](15), membershipclient.Pointer[int64](0))
    if err != nil {
        log.Fatal(err)
    }
    if res.ListApplicationsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `pageSize`                                                   | **int64*                                                     | :heavy_minus_sign:                                           | N/A                                                          |
| `page`                                                       | **int64*                                                     | :heavy_minus_sign:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ListOrganizationApplicationsResponse](../../pkg/models/operations/listorganizationapplicationsresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## GetOrganizationApplication

Get application for organization

### Example Usage

<!-- UsageSnippet language="go" operationID="getOrganizationApplication" method="get" path="/organizations/{organizationId}/applications/{applicationId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.GetOrganizationApplication(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.GetApplicationResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `applicationID`                                              | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.GetOrganizationApplicationResponse](../../pkg/models/operations/getorganizationapplicationresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## EnableApplicationForOrganization

Enable application for organization

### Example Usage

<!-- UsageSnippet language="go" operationID="enableApplicationForOrganization" method="put" path="/organizations/{organizationId}/applications/{applicationId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.EnableApplicationForOrganization(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.EnableApplicationForOrganizationResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `applicationID`                                              | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.EnableApplicationForOrganizationResponse](../../pkg/models/operations/enableapplicationfororganizationresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## DisableApplicationForOrganization

Disable application for organization

### Example Usage

<!-- UsageSnippet language="go" operationID="disableApplicationForOrganization" method="delete" path="/organizations/{organizationId}/applications/{applicationId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.DisableApplicationForOrganization(ctx, "<id>", "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `organizationID`                                             | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `applicationID`                                              | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.DisableApplicationForOrganizationResponse](../../pkg/models/operations/disableapplicationfororganizationresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ListApplications

List applications

### Example Usage

<!-- UsageSnippet language="go" operationID="listApplications" method="get" path="/applications" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ListApplications(ctx)
    if err != nil {
        log.Fatal(err)
    }
    if res.ListApplicationsResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ListApplicationsResponse](../../pkg/models/operations/listapplicationsresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## CreateApplication

Create application

### Example Usage

<!-- UsageSnippet language="go" operationID="createApplication" method="post" path="/applications" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.CreateApplication(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.CreateApplicationResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                            | Type                                                                 | Required                                                             | Description                                                          |
| -------------------------------------------------------------------- | -------------------------------------------------------------------- | -------------------------------------------------------------------- | -------------------------------------------------------------------- |
| `ctx`                                                                | [context.Context](https://pkg.go.dev/context#Context)                | :heavy_check_mark:                                                   | The context to use for the request.                                  |
| `request`                                                            | [shared.ApplicationData](../../pkg/models/shared/applicationdata.md) | :heavy_check_mark:                                                   | The request object to use for the request.                           |
| `opts`                                                               | [][operations.Option](../../pkg/models/operations/option.md)         | :heavy_minus_sign:                                                   | The options for this request.                                        |

### Response

**[*operations.CreateApplicationResponse](../../pkg/models/operations/createapplicationresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## GetApplication

Get application

### Example Usage

<!-- UsageSnippet language="go" operationID="getApplication" method="get" path="/applications/{applicationId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.GetApplication(ctx, "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.GetApplicationResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `applicationID`                                              | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.GetApplicationResponse](../../pkg/models/operations/getapplicationresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## UpdateApplication

Update application

### Example Usage

<!-- UsageSnippet language="go" operationID="updateApplication" method="put" path="/applications/{applicationId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.UpdateApplication(ctx, "<id>", nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.UpdateApplicationResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                             | Type                                                                  | Required                                                              | Description                                                           |
| --------------------------------------------------------------------- | --------------------------------------------------------------------- | --------------------------------------------------------------------- | --------------------------------------------------------------------- |
| `ctx`                                                                 | [context.Context](https://pkg.go.dev/context#Context)                 | :heavy_check_mark:                                                    | The context to use for the request.                                   |
| `applicationID`                                                       | *string*                                                              | :heavy_check_mark:                                                    | N/A                                                                   |
| `applicationData`                                                     | [*shared.ApplicationData](../../pkg/models/shared/applicationdata.md) | :heavy_minus_sign:                                                    | N/A                                                                   |
| `opts`                                                                | [][operations.Option](../../pkg/models/operations/option.md)          | :heavy_minus_sign:                                                    | The options for this request.                                         |

### Response

**[*operations.UpdateApplicationResponse](../../pkg/models/operations/updateapplicationresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## DeleteApplication

Delete application

### Example Usage

<!-- UsageSnippet language="go" operationID="deleteApplication" method="delete" path="/applications/{applicationId}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.DeleteApplication(ctx, "<id>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `applicationID`                                              | *string*                                                     | :heavy_check_mark:                                           | N/A                                                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.DeleteApplicationResponse](../../pkg/models/operations/deleteapplicationresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## CreateApplicationScope

Create application scope

### Example Usage

<!-- UsageSnippet language="go" operationID="createApplicationScope" method="post" path="/applications/{applicationId}/scopes" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.CreateApplicationScope(ctx, "550e8400-e29b-41d4-a716-446655440000", nil)
    if err != nil {
        log.Fatal(err)
    }
    if res.CreateApplicationScopeResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                                         | Type                                                                                              | Required                                                                                          | Description                                                                                       | Example                                                                                           |
| ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| `ctx`                                                                                             | [context.Context](https://pkg.go.dev/context#Context)                                             | :heavy_check_mark:                                                                                | The context to use for the request.                                                               |                                                                                                   |
| `applicationID`                                                                                   | *string*                                                                                          | :heavy_check_mark:                                                                                | The unique identifier of the application (UUID format)                                            | 550e8400-e29b-41d4-a716-446655440000                                                              |
| `createApplicationScopeRequest`                                                                   | [*shared.CreateApplicationScopeRequest](../../pkg/models/shared/createapplicationscoperequest.md) | :heavy_minus_sign:                                                                                | N/A                                                                                               |                                                                                                   |
| `opts`                                                                                            | [][operations.Option](../../pkg/models/operations/option.md)                                      | :heavy_minus_sign:                                                                                | The options for this request.                                                                     |                                                                                                   |

### Response

**[*operations.CreateApplicationScopeResponse](../../pkg/models/operations/createapplicationscoperesponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## DeleteApplicationScope

Delete a specific scope from an application. This operation requires system administrator privileges.

### Example Usage

<!-- UsageSnippet language="go" operationID="deleteApplicationScope" method="delete" path="/applications/{applicationId}/scopes/{scopeID}" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.DeleteApplicationScope(ctx, "550e8400-e29b-41d4-a716-446655440000", 115177)
    if err != nil {
        log.Fatal(err)
    }
    if res.Error != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  | Example                                                      |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |                                                              |
| `applicationID`                                              | *string*                                                     | :heavy_check_mark:                                           | The unique identifier of the application (UUID format)       | 550e8400-e29b-41d4-a716-446655440000                         |
| `scopeID`                                                    | *int64*                                                      | :heavy_check_mark:                                           | The unique identifier of the scope to operate on             |                                                              |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |                                                              |

### Response

**[*operations.DeleteApplicationScopeResponse](../../pkg/models/operations/deleteapplicationscoperesponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.Error    | 400, 404           | application/json   |
| sdkerrors.Error    | 500                | application/json   |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## CreateUser

Create a new user in the system. This operation requires system administrator privileges.

### Example Usage

<!-- UsageSnippet language="go" operationID="createUser" method="post" path="/users" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.CreateUser(ctx, shared.CreateUserRequest{
        Email: "user@example.com",
    })
    if err != nil {
        log.Fatal(err)
    }
    if res.CreateUserResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                                | Type                                                                     | Required                                                                 | Description                                                              |
| ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ | ------------------------------------------------------------------------ |
| `ctx`                                                                    | [context.Context](https://pkg.go.dev/context#Context)                    | :heavy_check_mark:                                                       | The context to use for the request.                                      |
| `request`                                                                | [shared.CreateUserRequest](../../pkg/models/shared/createuserrequest.md) | :heavy_check_mark:                                                       | The request object to use for the request.                               |
| `opts`                                                                   | [][operations.Option](../../pkg/models/operations/option.md)             | :heavy_minus_sign:                                                       | The options for this request.                                            |

### Response

**[*operations.CreateUserResponse](../../pkg/models/operations/createuserresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.Error    | 400                | application/json   |
| sdkerrors.Error    | 500                | application/json   |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |

## ReadConnectedUser

Read user

### Example Usage

<!-- UsageSnippet language="go" operationID="readConnectedUser" method="get" path="/me" -->
```go
package main

import(
	"context"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"log"
)

func main() {
    ctx := context.Background()

    s := membershipclient.New(
        membershipclient.WithSecurity(shared.Security{
            Oauth2: "<YOUR_OAUTH2_HERE>",
        }),
    )

    res, err := s.ReadConnectedUser(ctx)
    if err != nil {
        log.Fatal(err)
    }
    if res.ReadUserResponse != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                    | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `ctx`                                                        | [context.Context](https://pkg.go.dev/context#Context)        | :heavy_check_mark:                                           | The context to use for the request.                          |
| `opts`                                                       | [][operations.Option](../../pkg/models/operations/option.md) | :heavy_minus_sign:                                           | The options for this request.                                |

### Response

**[*operations.ReadConnectedUserResponse](../../pkg/models/operations/readconnecteduserresponse.md), error**

### Errors

| Error Type         | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4XX, 5XX           | \*/\*              |