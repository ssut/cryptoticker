package cryptoticker

import "github.com/parnurzeal/gorequest"

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

type bithumbDataMap = map[string]bithumbMarketTicker

type bithumbTicker struct {
	Status string
	Data   bithumbDataMap
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
	pairs := []*CurrencyPair{}
	for key := range t.Data {
		pair := &CurrencyPair{"KRW", key}
		pairs = append(pairs, pair)
	}
	return pairs, nil
}

func (t *bithumbTicker) Tickers() ([]*Ticker, error) {
	tickers := []*Ticker{}
	coins, _ := t.Coins()

	for _, currency := range coins {
		data := t.Data[currency.Next]
		ticker := &Ticker{
			currency,
			data.Volume1D,
			data.Last,
			data.High,
			data.Low,
			data.First,
		}
		tickers = append(tickers, ticker)
	}

	return tickers, nil
}
