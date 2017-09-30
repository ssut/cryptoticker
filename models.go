package cryptoticker

// TickerType represents makerts
type TickerType int

const (
	// Coinone
	CoinoneTicker  TickerType = 1 << iota
	BithumbTicker             // Bithumb
	PoloniexTicker            // Poloniex
	BittrexTicker             // Bittrex
	BitfinexTicker            // Bitfinex
)

// SubTickerType represents a subscribable markets
type SubTickerType int

const (
	// Coinone
	CoinoneSubTicker  SubTickerType = 1 << iota
	PoloniexSubTicker               // Poloniex
)
