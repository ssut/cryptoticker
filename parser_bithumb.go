package cryptoticker

import "github.com/parnurzeal/gorequest"

const (
	BithumbParserBaseURL    = "https://api.bithumb.com"
	BithumbParserTickerPath = "/public/ticker/ALL"
)

type BithumbMarketTicker struct {
	Volume1D string `json:"volume_1day"`
	Volume7D string `json:"volume_7day"`
	Last     string `json:"closing_price"`
	High     string `json:"max_price"`
	Low      string `json:"min_price"`
	First    string `json:"opening_price"`
}

type BithumbDataMap = map[string]BithumbMarketTicker

type BithumbTicker struct {
	Status string
	Data   BithumbDataMap
}

type BithumbParser struct {
	client *gorequest.SuperAgent
}

func NewBithumbParser() *BithumbParser {
	parser := &BithumbParser{
		client: gorequest.New(),
	}
	return parser
}

func (p *BithumbParser) TickerURLString() string {
	return BithumbParserBaseURL + BithumbParserTickerPath
}

func (p *BithumbParser) RawTicker() (IParsableTicker, error) {
	var ticker BithumbTicker
	_, _, err := p.client.Get(p.TickerURLString()).EndStruct(&ticker)
	if err != nil {
		return nil, err[0]
	}

	return &ticker, nil
}

func (t *BithumbTicker) Coins() ([]*CurrencyPair, error) {
	pairs := []*CurrencyPair{}
	for key := range t.Data {
		pair := &CurrencyPair{"KRW", key}
		pairs = append(pairs, pair)
	}
	return pairs, nil
}

func (t *BithumbTicker) Tickers() ([]*ParserTicker, error) {
	tickers := []*ParserTicker{}
	coins, _ := t.Coins()

	for _, currency := range coins {
		data := t.Data[currency.Next]
		ticker := &ParserTicker{
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
