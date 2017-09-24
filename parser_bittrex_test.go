package cryptoticker

import (
	"testing"

	"github.com/parnurzeal/gorequest"
	. "github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

var (
	bittrexTestParser *bittrexParser
)

func prepareBittrexParser() {
	bittrexTestParser = newBittrexParser()
	gorequest.DisableTransportSwap = true
}

func TestBittrexParser(t *testing.T) {
	prepareBittrexParser()
	defer gock.Off()

	Convey("Given a sample data", t, func() {
		raw := map[string]interface{}{
			"success": true,
			"message": "",
			"result": []map[string]interface{}{
				{
					"MarketName": "BTC-ETH",
					"High":       0.12345678,
					"Low":        0.12345678,
					"Volume":     10000.12345678,
					"Last":       0.12345678,
					"PrevDay":    0.12345678,
				},
			},
		}

		gock.New(bittrexParserBaseURL).
			Get(bittrexParserTickerPath).
			Reply(200).
			JSON(raw)

		Convey(".RawTicker() should return a BittrexTicker that contains the same data", func() {
			expected := &bittrexTicker{
				Success: true,
				Message: "",
				Result: []*bittrexMarketTicker{
					{
						Currency: "BTC-ETH",
						High:     0.12345678,
						Low:      0.12345678,
						Volume:   10000.12345678,
						Last:     0.12345678,
						First:    0.12345678,
					},
				},
			}
			actual, err := bittrexTestParser.RawTicker()
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)

			Convey(".Ticker() should return a list of BittrexParser", func() {
				expected := []*Ticker{
					{
						&CurrencyPair{"BTC", "ETH"},
						"10000.12345678",
						"0.12345678",
						"0.12345678",
						"0.12345678",
						"0.12345678",
					},
				}
				tickers, err := actual.Tickers()
				So(err, ShouldBeNil)
				So(tickers, ShouldResemble, expected)
			})
		})
	})
}
