package cryptoticker

import "strings"

// TickerType represents makerts
type TickerType int

const (
	CoinoneTicker  TickerType = 1 << iota // Coinone
	BithumbTicker                         // Bithumb
	PoloniexTicker                        // Poloniex
	BittrexTicker                         // Bittrex
	BitfinexTicker                        // Bitfinex
)

// IParsableTicker represents a object that can parse and return Ticker
type IParsableTicker interface {
	Coins() ([]*CurrencyPair, error)
	Tickers() ([]*Ticker, error)
}

// IParser represents a object that can parse and return a IParsableTicker
type IParser interface {
	// Coins returns an array of supported cryptocurrencies
	RawTicker() (IParsableTicker, error)
}

// Ticker represents a Ticker
type Ticker struct {
	Currency *CurrencyPair
	Volume   string
	Last     string
	High     string
	Low      string
	First    string
}

// CurrencyPair represents a pair of trading currency
type CurrencyPair struct {
	Base string
	Next string
}

// String returns a unique string based on the given pair
func (p *CurrencyPair) String() string {
	return strings.Join([]string{p.Base, p.Next}, "_")
}

// Parser represents a Crypto Client
type Parser struct {
	parser IParser
}

// NewParser returns a new Crypto Client with specific market
func NewParser(p TickerType) *Parser {
	var parser IParser
	switch p {
	case CoinoneTicker:
		parser = newCoinoneParser()
	case BithumbTicker:
		parser = newBithumbParser()
	case PoloniexTicker:
		parser = newPoloniexParser()
	case BittrexTicker:
		parser = newBittrexParser()
	case BitfinexTicker:
		parser = newBitfinexParser()
	}

	return &Parser{parser: parser}
}

// Parse tries to fetch and parse ticker data from the market, and returns IParsableTicker
func (p *Parser) Parse() (IParsableTicker, error) {
	ticker, err := p.parser.RawTicker()
	return ticker, err
}
