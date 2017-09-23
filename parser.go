package cryptoticker

import "strings"

// ParserType represents makerts
type ParserType int

const (
	CoinoneParserType  ParserType = 1 << iota // Coinone
	BithumbParserType                         // Bithumb
	PoloniexParserType                        // Poloniex
	BittrexParerType                          // Bittrex
)

// IParsableTicker represents a object that can parse and return ParserTicker
type IParsableTicker interface {
	Coins() ([]*CurrencyPair, error)
	Tickers() ([]*ParserTicker, error)
}

// IParser represents a object that can parse and return a IParsableTicker
type IParser interface {
	// Coins returns an array of supported cryptocurrencies
	RawTicker() (IParsableTicker, error)
}

// ParserTicker represents a Ticker
type ParserTicker struct {
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
func NewParser(p ParserType) *Parser {
	var parser IParser
	switch p {
	case CoinoneParserType:
		parser = NewCoinoneParser()
	case BithumbParserType:
		parser = NewBithumbParser()
	case PoloniexParserType:
		parser = NewPoloniexParser()
	case BittrexParerType:
		parser = NewBittrexParser()
	}

	return &Parser{parser: parser}
}
