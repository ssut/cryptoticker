package cryptoticker

import (
	"strings"

	"github.com/parnurzeal/gorequest"
	"github.com/shopspring/decimal"
)

const (
	poloniexParserBaseURL    = "https://www.poloniex.com"
	poloniexParserTickerPath = "/public?command=returnTicker"
)

type poloniexMarketTicker struct {
	ID          int
	Volume      string `json:"quoteVolume"`
	Last        string
	High        string `json:"high24hr"`
	Low         string `json:"low24hr"`
	PercentDiff string `json:"percentChange"`
}

type poloniexTicker map[string]poloniexMarketTicker

type poloniexParser struct {
	client *gorequest.SuperAgent
}

func newPoloniexParser() *poloniexParser {
	parser := &poloniexParser{
		client: gorequest.New(),
	}
	return parser
}

func (p *poloniexParser) TickerURLString() string {
	return poloniexParserBaseURL + poloniexParserTickerPath
}

func (p *poloniexParser) RawTicker() (IParsableTicker, error) {
	var ticker poloniexTicker
	_, _, err := p.client.Get(p.TickerURLString()).EndStruct(&ticker)
	if err != nil {
		return nil, err[0]
	}

	return &ticker, nil
}

func (t *poloniexTicker) Coins() ([]*CurrencyPair, error) {
	pairs := make([]*CurrencyPair, len(*t))
	var i int
	for key := range *t {
		sep := strings.Split(key, "_")
		pair := &CurrencyPair{sep[0], sep[1]}
		pairs[i] = pair
		i++
	}

	return pairs, nil
}

func (t *poloniexTicker) Tickers() ([]*Ticker, error) {
	tickers := []*Ticker{}
	coins, _ := t.Coins()

	for _, currency := range coins {
		data := (*t)[currency.String()]
		lastDec, err := decimal.NewFromString(data.Last)
		if err != nil {
			return nil, err
		}
		diffDec, err := decimal.NewFromString(data.PercentDiff)
		if err != nil {
			return nil, err
		}

		first := lastDec.Sub(lastDec.Mul(diffDec))
		ticker := &Ticker{
			currency,
			data.Volume,
			data.Last,
			data.High,
			data.Low,
			first.String(),
		}
		tickers = append(tickers, ticker)
	}

	return tickers, nil
}
