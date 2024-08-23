Kraken GO API Client
====================

Forked from [svanas/kraken-go-api-client](https://github.com/svanas/kraken-go-api-client)

A simple API Client for the [Kraken](https://www.kraken.com/ "Kraken") Trading platform.

Example usage:

```go
package main

import (
	"fmt"
	"log"

	"github.com/sergey-lipin/kraken-go-api-client"
)

func main() {
	api := krakenapi.New("KEY", "SECRET")
	result, err := api.Query("Ticker", map[string]string{
		"pair": "XXBTZEUR",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Result: %+v\n", result)

	// There are also some strongly typed methods available
	ticker, err := api.Ticker(krakenapi.XXBTZEUR)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ticker.XXBTZEUR.OpeningPrice)
}
```

## Contributors
 - Piega
 - Glavic
 - MarinX
 - bjorand
 - [khezen](https://github.com/khezen)
 