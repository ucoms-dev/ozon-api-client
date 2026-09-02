# Ozon Seller API Client
A Ozon Seller API client written in Go, maintained by UCOMS.

![example workflow](https://github.com/ucoms-dev/ozon-api-client/actions/workflows/tests.yml/badge.svg)

[Ozon](https://ozon.ru) is a marketplace for small and medium enterprises to launch and grow their businesses in Russia.

Read full [documentation](https://docs.ozon.ru/api/seller/en/#tag/Introduction)

## How to start
### API
Get Client-Id and Api-Key in your seller profile [here](https://seller.ozon.ru/app/settings/api-keys?locale=en)

Just add dependency to your project and you're ready to go.
```bash
go get github.com/ucoms-dev/ozon-api-client
```
A simple example on how to use this library:
```Golang
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ucoms-dev/ozon-api-client/ozon"
)

func main() {
	// Create a client with your Client-Id and Api-Key
	// [Documentation]: https://docs.ozon.ru/api/seller/en/#tag/Auth
	opts := []ozon.ClientOption{
		ozon.WithAPIKey("api-key"),
		ozon.WithClientId("client-id"),
	}
	c := ozon.NewClient(opts...)

	// Send request with parameters
	resp, err := c.Products().GetProductDetails(context.Background(), &ozon.GetProductDetailsParams{
		ProductId: 123456789,
	})
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalf("error when getting product details: %s", err)
	}

	// Do some stuff
	for _, d := range resp.Result.Barcodes {
		fmt.Printf("Barcode %s\n", d)
	}
}
```

### Notifications
Ozon can send push-notifications to your REST server. There is an implementation of REST server that handles notifications in this library.

[Official documentation](https://docs.ozon.ru/api/seller/en/#tag/push_intro)

How to use:
```Golang
package main

import (
	"log"

	"github.com/ucoms-dev/ozon-api-client/ozon/notifications"
)

func main() {
	// Create server
	port := 5000
	server := notifications.NewNotificationServer(port)

	// Register handlers passing message type and handler itself
	server.Register(notifications.ChatClosedType, func(req interface{}) error {
		notification := req.(*notifications.ChatClosed)

		// Do something with the notification here...
		log.Printf("chat %s has been closed\n", notification.ChatId)

		return nil
	})

	// Run server
	if err := server.Run(); err != nil {
		log.Printf("error while running notification server: %s", err)
	}
}
```

## API contract audit

The repository includes a deterministic audit tool that compares the HTTP
method/path pairs used by the Go client with an Ozon OpenAPI JSON document:

```bash
go run ./cmd/contract-audit \
  -swagger /path/to/swagger.json \
  -client ozon \
  -date YYYY-MM-DD \
  -output docs/ozon-seller-api-contract-audit-YYYY-MM-DD.md
```

The generated report records the source document SHA-256, exact matches,
client-only endpoints, method mismatches, deprecated endpoints, and operations
that are not implemented by the client. The current audited snapshot is in
[`docs/ozon-seller-api-contract-audit-2026-09-02.md`](docs/ozon-seller-api-contract-audit-2026-09-02.md).
The reviewed compatibility decisions and implementation order are in
[`docs/ozon-seller-api-contract-plan-2026-09-02.md`](docs/ozon-seller-api-contract-plan-2026-09-02.md).

Deprecated exported methods remain available throughout the v1 release line for
source compatibility. New code should use the replacement named in each Go doc
comment. Removing those methods is reserved for a future major version.
