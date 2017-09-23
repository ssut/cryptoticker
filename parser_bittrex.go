package cryptoticker

import (
	"strings"

	"github.com/parnurzeal/gorequest"
	"github.com/shopspring/decimal"
)

const (
	BittrexParserBaseURL    = "https://bittrex.com"
	BittrexParserTickerPath = "/api/v1.1/public/getmarketsummaries"
)

type BittrexMarketTicker struct {
	Currency string `json:"MarketName"`
	Volume   float64
	Last     float64
	High     float64
	Low      float64
	First    float64 `json:"PrevDay"`
}

type BittrexTicker struct {
	Success bool
	Message string
	Result  []*BittrexMarketTicker
}

type BittrexParser struct {
	client *gorequest.SuperAgent
}

func NewBittrexParser() *BittrexParser {
	parser := &BittrexParser{
		client: gorequest.New(),
	}
	return parser
}

func (p *BittrexParser) TickerURLString() string {
	return BittrexParserBaseURL + BittrexParserTickerPath
}

func (p *BittrexParser) RawTicker() (IParsableTicker, error) {
	var ticker BittrexTicker
	_, _, err := p.client.Get(p.TickerURLString()).EndStruct(&ticker)
	if err != nil {
		return nil, err[0]
	}

	return &ticker, nil
}

func (t *BittrexTicker) Coins() ([]*CurrencyPair, error) {
	pairs := make([]*CurrencyPair, len(t.Result))
	for i, coin := range t.Result {
		sep := strings.Split(coin.Currency, "-")
		pair := &CurrencyPair{sep[0], sep[1]}
		pairs[i] = pair
	}

	return pairs, nil
}

func (t *BittrexTicker) Tickers() ([]*ParserTicker, error) {
	tickers := make([]*ParserTicker, len(t.Result))
	coins, _ := t.Coins()

	for i, currency := range coins {
		data := t.Result[i]
		ticker := &ParserTicker{
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
