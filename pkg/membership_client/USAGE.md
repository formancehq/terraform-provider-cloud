<!-- Start SDK Example Usage [usage] -->
```go
package main

import (
	"context"
	formancesdkcloudgo "github.com/formancehq/formance-sdk-cloud-go"
	"github.com/formancehq/formance-sdk-cloud-go/pkg/models/shared"
	"log"
)

func main() {
	ctx := context.Background()

	s := formancesdkcloudgo.New(
		formancesdkcloudgo.WithSecurity(shared.Security{
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