package cryptoticker

import (
	"strings"

	"fmt"

	"reflect"

	"github.com/fatih/structs"
	"github.com/parnurzeal/gorequest"
	"github.com/shopspring/decimal"
)

const (
	bitfinexParserBaseURL     = "https://api.bitfinex.com"
	bitfinexParserSymbolsPath = "/v1/symbols"
	bitfinexParserTickersPath = "/v2/tickers"
)

var (
	bitfinexMarketTickerOrder = [...]string{
		"Currency",
		"Bid",
		"BidSize",
		"Ask",
		"AskSize",
		"Diff",
		"DiffPercent",
		"Last",
		"Volume",
		"High",
		"Low",
	}
)

type bitfinexMarketTickerProp struct {
	Name string
	Type reflect.Kind
}

type bitfinexMarketTicker struct {
	Currency    *CurrencyPair
	Bid         float64
	BidSize     float64
	Ask         float64
	AskSize     float64
	Diff        float64
	DiffPercent float64
	Last        float64
	Volume      float64
	High        float64
	Low         float64
}

type bitfinexTicker struct {
	pairs   []*CurrencyPair
	tickers []*bitfinexMarketTicker
}

type bitfinexParser struct {
	client *gorequest.SuperAgent
}

func newBitfinexParser() *bitfinexParser {
	parser := &bitfinexParser{
		client: gorequest.New(),
	}
	return parser
}

func (p *bitfinexParser) TickerURLString() string {
	return bitfinexParserBaseURL + bitfinexParserTickersPath
}

func (p *bitfinexParser) SymbolsURLString() string {
	return bitfinexParserBaseURL + bitfinexParserSymbolsPath
}

func (p *bitfinexParser) symbols() ([]*CurrencyPair, error) {
	var symbols []string
	_, _, err := p.client.Get(p.SymbolsURLString()).EndStruct(&symbols)
	if err != nil {
		return nil, err[0]
	}

	pairs := []*CurrencyPair{}
	for _, symbol := range symbols {
		base := strings.ToUpper(symbol[:3])
		next := strings.ToUpper(symbol[3:])
		pair := &CurrencyPair{base, next}
		pairs = append(pairs, pair)
	}

	return pairs, nil
}

func (p *bitfinexParser) RawTicker() (IParsableTicker, error) {
	pairs, err := p.symbols()
	if err != nil {
		return nil, err
	}

	symbols := make([]string, len(pairs))
	for i, pair := range pairs {
		symbols[i] = fmt.Sprintf("t%s%s", pair.Base, pair.Next)
	}
	symbolsParam := strings.Join(symbols, ",")

	var items []interface{}
	_, _, errs := p.client.Get(p.TickerURLString()).
		Param("symbols", symbolsParam).
		EndStruct(&items)
	if errs != nil {
		return nil, errs[0]
	}

	// Unpack items
	t := &bitfinexTicker{
		pairs:   pairs,
		tickers: []*bitfinexMarketTicker{},
	}
	for _, raw := range items {
		market := &bitfinexMarketTicker{}
		s := structs.New(market)
		for i, item := range raw.([]interface{}) {
			// possible field name
			fieldName := bitfinexMarketTickerOrder[i]
			if item, ok := item.(string); ok {
				for _, p := range pairs {
					if fmt.Sprintf("%s%s", p.Base, p.Next) == strings.ToUpper(item[1:]) {
						market.Currency = p
					}
				}
			}
			if item, ok := item.(float64); ok {
				if field, ok := s.FieldOk(fieldName); ok {
					field.Set(item)
				}
			}
		}

		if market.Currency != nil {
			t.tickers = append(t.tickers, market)
		}
	}

	return t, nil
}

func (t *bitfinexTicker) Coins() ([]*CurrencyPair, error) {
	return t.pairs, nil
}

func (t *bitfinexTicker) Tickers() ([]*Ticker, error) {
	tickers := make([]*Ticker, len(t.tickers))
	for i, ticker := range t.tickers {
		marketVolume := decimal.NewFromFloat(ticker.Volume)
		lastPrice := decimal.NewFromFloat(ticker.Last)
		highPrice := decimal.NewFromFloat(ticker.High)
		lowPrice := decimal.NewFromFloat(ticker.Low)
		firstPrice := lastPrice.Sub(decimal.NewFromFloat(ticker.Diff))
		tickers[i] = &Ticker{
			ticker.Currency,
			marketVolume.String(),
			lastPrice.String(),
			highPrice.String(),
			lowPrice.String(),
			firstPrice.String(),
		}
	}
	return tickers, nil
}
