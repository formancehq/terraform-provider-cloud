<!-- Start SDK Example Usage [usage] -->
```go
package main

import (
	"context"
	membershipclient "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
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
<!-- End SDK Example Usage [usage] -->