package cryptoticker

import (
	"testing"

	"github.com/parnurzeal/gorequest"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

var (
	bitfinexTestParser *bitfinexParser
)

func prepareBitfinexParser() {
	bitfinexTestParser = newBitfinexParser()
	gorequest.DisableTransportSwap = true
}

func TestBitfinexParser(t *testing.T) {
	prepareBitfinexParser()
	defer gock.Off()

	prepareMockServer := func() {
		raw := `["btcusd", "ethusd"]`
		rawTicker := `[["tBTCUSD",3876.1,110.74753748,3877.7,75.26554118,172.5,0.0466,3875.9,27032.68851049,3907.9,3638.1],["tETHUSD",293.3,401.86437684,293.48,404.07426443,8.29,0.0291,293.3,169154.90556116,295.44,279.07]]`
		gock.New(bitfinexParserBaseURL).
			Get(bitfinexParserSymbolsPath).
			Reply(200).
			SetHeader("Content-Type", "application/json").
			BodyString(raw)
		gock.New(bitfinexParserBaseURL).
			Get(bitfinexParserTickersPath).
			MatchParam("symbols", "tBTCUSD,tETHUSD").
			Reply(200).
			SetHeader("Content-Type", "application/json").
			BodyString(rawTicker)
	}

	Convey("Given a sample symbols and tickers", t, func() {
		prepareMockServer()

		Convey(".symbols() should return a list of supported symbols in CurrencyPair", func() {
			pairs := []*CurrencyPair{
				{"BTC", "USD"},
				{"ETH", "USD"},
			}
			symbols, err := bitfinexTestParser.symbols()
			So(err, ShouldBeNil)
			So(symbols, ShouldResemble, pairs)

			Convey(".RawTicker() should return a IParsableTicker", func() {
				prepareMockServer()
				expected := &bitfinexTicker{
					pairs: symbols,
					tickers: []*bitfinexMarketTicker{
						{
							pairs[0],
							3876.1,
							110.74753748,
							3877.7,
							75.26554118,
							172.5,
							0.0466,
							3875.9,
							27032.68851049,
							3907.9,
							3638.1,
						},
						{
							pairs[1],
							293.3,
							401.86437684,
							293.48,
							404.07426443,
							8.29,
							0.0291,
							293.3,
							169154.90556116,
							295.44,
							279.07,
						},
					},
				}
				ticker, err := bitfinexTestParser.RawTicker()
				So(err, ShouldBeNil)
				So(ticker.(*bitfinexTicker).pairs, ShouldResemble, expected.pairs)
				So(ticker.(*bitfinexTicker).tickers, ShouldResemble, expected.tickers)

				Convey(".Coins() should return the same pairs as .symbols()", func() {
					coins, err := ticker.Coins()
					So(err, ShouldBeNil)
					So(coins, ShouldResemble, symbols)
				})

				Convey(".Ticker() should return a list of Ticker", func() {
					expected := []*Ticker{
						{
							pairs[0],
							"27032.68851049",
							"3875.9",
							"3907.9",
							"3638.1",
							"3703.4",
						},
						{
							pairs[1],
							"169154.90556116",
							"293.3",
							"295.44",
							"279.07",
							"285.01",
						},
					}
					actual, err := ticker.Tickers()
					So(err, ShouldBeNil)
					So(actual, ShouldResemble, expected)
				})
			})
		})

	})
}
