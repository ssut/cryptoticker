package cryptoticker

import (
	"github.com/fatih/structs"
	"github.com/parnurzeal/gorequest"
)

const (
	CoinoneParserBaseURL    = "https://api.coinone.co.kr"
	CoinoneParserTickerPath = "/ticker/?format=json&currency="
)

type CoinoneMarketTicker struct {
	Currency string
	Volume   string
	Last     string
	High     string
	Low      string
	First    string
}

type CoinoneTicker struct {
	Result bool
	BTC    *CoinoneMarketTicker
	BCH    *CoinoneMarketTicker
	ETH    *CoinoneMarketTicker
	ETC    *CoinoneMarketTicker
	XRP    *CoinoneMarketTicker
	QTUM   *CoinoneMarketTicker
}

type CoinoneParser struct {
	client *gorequest.SuperAgent
}

func NewCoinoneParser() *CoinoneParser {
	parser := &CoinoneParser{
		client: gorequest.New(),
	}
	return parser
}

func (p *CoinoneParser) TickerURLString() string {
	return CoinoneParserBaseURL + CoinoneParserTickerPath
}

func (p *CoinoneParser) RawTicker() (IParsableTicker, error) {
	var ticker CoinoneTicker
	_, _, err := p.client.Get(p.TickerURLString()).EndStruct(&ticker)
	if err != nil {
		// return the first error because it contains multiple errors
		return nil, err[0]
	}

	return &ticker, nil
}

// Coins returns a list of supported cryptocurrency pairs
func (t *CoinoneTicker) Coins() ([]*CurrencyPair, error) {
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

// Tickers returns a list of ParserTicker
func (t *CoinoneTicker) Tickers() ([]*ParserTicker, error) {
	tickers := []*ParserTicker{}
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
				ct := t.(*CoinoneMarketTicker)

				ticker := &ParserTicker{
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
