package cryptoticker

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCurrencyPair_String(t *testing.T) {
	Convey("Given a sample data", t, func() {
		samples := []*CurrencyPair{
			{"BTC", "ETH"},
			{"KRW", "BTC"},
			{"USDT", "BTC"},
		}
		expected := []string{"BTC_ETH", "KRW_BTC", "USDT_BTC"}

		Convey(".String() should return uniquely identifiable name", func() {
			for i, pair := range samples {
				So(pair.String(), ShouldEqual, expected[i])
			}
		})
	})
}

func TestNewParser(t *testing.T) {

}