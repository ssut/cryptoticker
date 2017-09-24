package cryptoticker

import (
	"github.com/fatih/structs"
	"github.com/parnurzeal/gorequest"
)

const (
	coinoneParserBaseURL    = "https://api.coinone.co.kr"
	coinoneParserTickerPath = "/ticker/?format=json&currency="
)

type coinoneMarketTicker struct {
	Currency string
	Volume   string
	Last     string
	High     string
	Low      string
	First    string
}

type coinoneTicker struct {
	Result string
	BTC    *coinoneMarketTicker
	BCH    *coinoneMarketTicker
	ETH    *coinoneMarketTicker
	ETC    *coinoneMarketTicker
	XRP    *coinoneMarketTicker
	QTUM   *coinoneMarketTicker
}

type coinoneParser struct {
	client *gorequest.SuperAgent
}

func newCoinoneParser() *coinoneParser {
	parser := &coinoneParser{
		client: gorequest.New(),
	}
	return parser
}

func (p *coinoneParser) TickerURLString() string {
	return coinoneParserBaseURL + coinoneParserTickerPath
}

func (p *coinoneParser) RawTicker() (IParsableTicker, error) {
	var ticker coinoneTicker
	_, _, err := p.client.Get(p.TickerURLString()).EndStruct(&ticker)
	if err != nil {
		// return the first error because it contains multiple errors
		return nil, err[0]
	}

	return &ticker, nil
}

// Coins returns a list of supported cryptocurrency pairs
func (t *coinoneTicker) Coins() ([]*CurrencyPair, error) {
	pairs := []*CurrencyPair{
		{"KRW", "BTC"},
		{"KRW", "BCH"},
		{"KRW", "ETH"},
		{"KRW", "ETC"},
		{"KRW", "XRP"},
		{"KRW", "QTUM"},
	}
	return pairs, nil
}

// Tickers returns a list of Ticker
func (t *coinoneTicker) Tickers() ([]*Ticker, error) {
	tickers := []*Ticker{}
	coins, _ := t.Coins()

	s := structs.New(t)
	fields := s.Fields()
	for _, field := range fields {
		name := field.Name()
		for _, currency := range coins {
			if currency.Next == name {
				t := field.Value()
				if t == nil || field.IsZero() {
					break
				}
				ct := t.(*coinoneMarketTicker)

				ticker := &Ticker{
					currency,
					ct.Volume,
					ct.Last,
					ct.High,
					ct.Low,
					ct.First,
				}
				tickers = append(tickers, ticker)
				break
			}
		}
	}

	return tickers, nil
}
