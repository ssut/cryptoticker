package cryptoticker

import (
	"github.com/fatih/structs"
	"github.com/parnurzeal/gorequest"
)

const (
	bithumbParserBaseURL    = "https://api.bithumb.com"
	bithumbParserTickerPath = "/public/ticker/ALL"
)

type bithumbMarketTicker struct {
	Volume1D string `json:"volume_1day"`
	Volume7D string `json:"volume_7day"`
	Last     string `json:"closing_price"`
	High     string `json:"max_price"`
	Low      string `json:"min_price"`
	First    string `json:"opening_price"`
}

type bithumbTicker struct {
	Status string
	Data   struct {
		BTC  *bithumbMarketTicker
		ETH  *bithumbMarketTicker
		DASH *bithumbMarketTicker
		LTC  *bithumbMarketTicker
		ETC  *bithumbMarketTicker
		XRP  *bithumbMarketTicker
		BCH  *bithumbMarketTicker
		XMR  *bithumbMarketTicker
	}
}

type bithumbParser struct {
	client *gorequest.SuperAgent
}

func newBithumbParser() *bithumbParser {
	parser := &bithumbParser{
		client: gorequest.New(),
	}
	return parser
}

func (p *bithumbParser) TickerURLString() string {
	return bithumbParserBaseURL + bithumbParserTickerPath
}

func (p *bithumbParser) RawTicker() (IParsableTicker, error) {
	var ticker bithumbTicker
	_, _, err := p.client.Get(p.TickerURLString()).EndStruct(&ticker)
	if err != nil {
		return nil, err[0]
	}

	return &ticker, nil
}

func (t *bithumbTicker) Coins() ([]*CurrencyPair, error) {
	pairs := []*CurrencyPair{
		&CurrencyPair{"KRW", "BTC"},
		&CurrencyPair{"KRW", "ETH"},
		&CurrencyPair{"KRW", "DASH"},
		&CurrencyPair{"KRW", "LTC"},
		&CurrencyPair{"KRW", "ETC"},
		&CurrencyPair{"KRW", "XRP"},
		&CurrencyPair{"KRW", "BCH"},
		&CurrencyPair{"KRW", "XMR"},
	}
	return pairs, nil
}

func (t *bithumbTicker) Tickers() ([]*Ticker, error) {
	tickers := []*Ticker{}
	coins, _ := t.Coins()

	s := structs.New(t.Data)
	fields := s.Fields()
	for _, field := range fields {
		name := field.Name()
		for _, currency := range coins {
			if currency.Next == name {
				t := field.Value()
				if t == nil || field.IsZero() {
					break
				}
				data := t.(*bithumbMarketTicker)
				ticker := &Ticker{
					currency,
					data.Volume1D,
					data.Last,
					data.High,
					data.Low,
					data.First,
				}
				tickers = append(tickers, ticker)
				break
			}
		}
	}

	return tickers, nil
}
