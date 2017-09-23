package cryptoticker

import (
	"strings"

	"github.com/parnurzeal/gorequest"
	"github.com/shopspring/decimal"
)

const (
	PoloniexParserBaseURL    = "https://www.poloniex.com"
	PoloniexParserTickerPath = "/public?command=returnTicker"
)

type PoloniexMarketTicker struct {
	ID          int
	Volume      string `json:"quoteVolume"`
	Last        string
	High        string `json:"high24hr"`
	Low         string `json:"low24hr"`
	PercentDiff string `json:"percentChange"`
}

type PoloniexTicker map[string]PoloniexMarketTicker

type PoloniexParser struct {
	client *gorequest.SuperAgent
}

func NewPoloniexParser() *PoloniexParser {
	parser := &PoloniexParser{
		client: gorequest.New(),
	}
	return parser
}

func (p *PoloniexParser) TickerURLString() string {
	return PoloniexParserBaseURL + PoloniexParserTickerPath
}

func (p *PoloniexParser) RawTicker() (IParsableTicker, error) {
	var ticker PoloniexTicker
	_, _, err := p.client.Get(p.TickerURLString()).EndStruct(&ticker)
	if err != nil {
		return nil, err[0]
	}

	return &ticker, nil
}

func (t *PoloniexTicker) Coins() ([]*CurrencyPair, error) {
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

func (t *PoloniexTicker) Tickers() ([]*ParserTicker, error) {
	tickers := []*ParserTicker{}
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
		ticker := &ParserTicker{
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
