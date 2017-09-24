package cryptoticker

import (
	"testing"

	"github.com/parnurzeal/gorequest"
	. "github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

var (
	coinoneTestParser *coinoneParser
)

func prepareCoinoneParser() {
	coinoneTestParser = newCoinoneParser()
	gorequest.DisableTransportSwap = true
}

func TestCoinoneParser(t *testing.T) {
	prepareCoinoneParser()
	defer gock.Off()

	Convey("Given a sample data", t, func() {
		raw := map[string]interface{}{
			"result": true,
			"btc": map[string]string{
				"currency": "btc",
				"volume":   "100.0000",
				"last":     "1000000",
				"high":     "1000000",
				"low":      "1000000",
				"first":    "1000000",
			},
		}

		gock.New(coinoneParserBaseURL).
			Get("/ticker").
			MatchParam("format", "json").
			ParamPresent("currency").
			Reply(200).
			JSON(raw)

		Convey(".RawTicker() should return a CoinoneTicker that contains the same data", func() {
			expected := &coinoneTicker{
				Result: true,
				BTC:    &coinoneMarketTicker{"btc", "100.0000", "1000000", "1000000", "1000000", "1000000"},
			}
			actual, err := coinoneTestParser.RawTicker()
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)

			Convey(".Ticker() should return a list of Ticker", func() {
				expected := []*Ticker{
					{
						&CurrencyPair{"KRW", "BTC"},
						"100.0000",
						"1000000",
						"1000000",
						"1000000",
						"1000000",
					},
				}
				tickers, err := actual.Tickers()
				So(err, ShouldBeNil)
				So(tickers, ShouldResemble, expected)
			})
		})

	})

}
