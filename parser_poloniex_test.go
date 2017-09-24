package cryptoticker

import (
	"testing"

	"github.com/parnurzeal/gorequest"
	. "github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

var (
	poloniexTestParser *poloniexParser
)

func preparePoloParser() {
	poloniexTestParser = newPoloniexParser()
	gorequest.DisableTransportSwap = true
}

func TestPoloniexParser(t *testing.T) {
	preparePoloParser()
	defer gock.Off()

	Convey("Given a sample data", t, func() {
		raw := map[string]interface{}{
			"BTC_ETH": map[string]interface{}{
				"id":            1,
				"quoteVolume":   "1000.1234",
				"last":          "0.11111",
				"high24hr":      "0.22222",
				"low24hr":       "0.00001",
				"percentChange": "0.1",
			},
		}

		gock.New(poloniexParserBaseURL).
			Get("/public").
			MatchParam("command", "returnTicker").
			Reply(200).
			JSON(raw)

		Convey(".RawTicker() should return a poloniexTicker that contains the same data", func() {
			expected := &poloniexTicker{
				"BTC_ETH": {
					1,
					"1000.1234",
					"0.11111",
					"0.22222",
					"0.00001",
					"0.1",
				},
			}
			actual, err := poloniexTestParser.RawTicker()
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)

			Convey(".Ticker() should return a list of Ticker", func() {
				expected := []*Ticker{
					{
						&CurrencyPair{"BTC", "ETH"},
						"1000.1234",
						"0.11111",
						"0.22222",
						"0.00001",
						"0.099999",
					},
				}
				tickers, err := actual.Tickers()
				So(err, ShouldBeNil)
				So(tickers, ShouldResemble, expected)
			})
		})
	})
}
