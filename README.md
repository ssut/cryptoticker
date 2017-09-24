# CryptoTicker

Crypto-currency market ticker API provider around the biggest markets.



## Import

```go
import "github.com/ssut/cryptoticker"
```



## Usage

```go
package main

import (
	"fmt"
	"github.com/ssut/cryptoticker"
)

func main() {
    parser := cryptoticker.NewParser(cryptoticker.PoloniexTicker)
	parsed, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	coins, err := parsed.Coins()
	if err != nil {
		panic(err)
	}

	for _, coin := range coins {
		// coin.Base, coin.Next
		fmt.Println(coin.String())
	}

	tickers, err := parsed.Tickers()
	if err != nil {
		panic(err)
	}
	
  	for _, ticker := range tickers {
      // ticker.Currency (the same object as parsed.Coins())
      // ticker.Volume, ticker.Last, ticker.High, ticker.Low, ticker.First
  	}
}
```

