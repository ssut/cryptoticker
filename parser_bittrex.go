package cryptoticker

import (
	"strings"

	"github.com/parnurzeal/gorequest"
	"github.com/shopspring/decimal"
)

const (
	bittrexParserBaseURL    = "https://bittrex.com"
	bittrexParserTickerPath = "/api/v1.1/public/getmarketsummaries"
)

type bittrexMarketTicker struct {
	Currency string `json:"MarketName"`
	Volume   float64
	Last     float64
	High     float64
	Low      float64
	First    float64 `json:"PrevDay"`
}

type bittrexTicker struct {
	Success bool
	Message string
	Result  []*bittrexMarketTicker
}

type bittrexParser struct {
	client *gorequest.SuperAgent
}

func newBittrexParser() *bittrexParser {
	parser := &bittrexParser{
		client: gorequest.New(),
	}
	return parser
}

func (p *bittrexParser) TickerURLString() string {
	return bittrexParserBaseURL + bittrexParserTickerPath
}

func (p *bittrexParser) RawTicker() (IParsableTicker, error) {
	var ticker bittrexTicker
	_, _, err := p.client.Get(p.TickerURLString()).EndStruct(&ticker)
	if err != nil {
		return nil, err[0]
	}

	return &ticker, nil
}

func (t *bittrexTicker) Coins() ([]*CurrencyPair, error) {
	pairs := make([]*CurrencyPair, len(t.Result))
	for i, coin := range t.Result {
		sep := strings.Split(coin.Currency, "-")
		pair := &CurrencyPair{sep[0], sep[1]}
		pairs[i] = pair
	}

	return pairs, nil
}

func (t *bittrexTicker) Tickers() ([]*Ticker, error) {
	tickers := make([]*Ticker, len(t.Result))
	coins, _ := t.Coins()

	for i, currency := range coins {
		data := t.Result[i]
		ticker := &Ticker{
			currency,
			decimal.NewFromFloat(data.Volume).String(),
			decimal.NewFromFloat(data.Last).String(),
			decimal.NewFromFloat(data.High).String(),
			decimal.NewFromFloat(data.Low).String(),
			decimal.NewFromFloat(data.First).String(),
		}
		tickers[i] = ticker
	}

	return tickers, nil
}
