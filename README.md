# CryptoTicker

[![Build Status](https://travis-ci.org/ssut/cryptoticker.svg?branch=master)](https://travis-ci.org/ssut/cryptoticker)
[![Coverage Status](https://coveralls.io/repos/github/ssut/cryptoticker/badge.svg?branch=master)](https://coveralls.io/github/ssut/cryptoticker?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/ssut/cryptoticker)](https://goreportcard.com/report/github.com/ssut/cryptoticker)
[![GoDoc](https://godoc.org/github.com/ssut/cryptoticker?status.svg)](https://godoc.org/github.com/ssut/cryptoticker)


In short: **All-in-one cryptocurrency ticker pack**

Crypto-currency market ticker API provider around the biggest *worldwide* markets.

- [x] Bittrex (`BittrexTicker`)
- [x] Poloniex (`PoloniexTicker`)
- [x] Bitfinex (`Bitfinex`)
- [ ] GDAX
- [ ] Kraken
- [x] Coinone (`Coinone`)
- [x] Bithumb (`Bithumb`)
- [ ] bitFlyer

... Request for more markets are welcome!

## Quick Start

```go
go get -u github.com/ssut/cryptoticker
```

**example.go**

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

## Documentation

Note that this library is unfinished. Because of that there may be major changes to library in the future.

See the document at:

* [![GoDoc](https://godoc.org/github.com/ssut/cryptoticker?status.svg)](https://godoc.org/github.com/ssut/cryptoticker)
* Hand crafted documentation coming eventually.

## License

CryptoTicker is free software licensed under the MIT license. Details provided in the LICENSE file.

## Buy me a coffee?

If you feel like buying me a coffee, donations are welcome to:

```
BTC: 1PkTzRLgcggpfrj4UDqM1PksPYirJs8WqY
ETH: 0x159D5BDCD4971E8CF8F51C5FA01E8B1A8CD25FE1
XMR: 43JCpUNQuH9STRr6nuBWWXULC8qRFRwWhiMd98yJLinVWDiqKZPenszX3B76GD7fv5dBN5uXY78CLP3pknQh9HhyR7ohVCU
```

