package cryptoticker

import (
	"time"

	"github.com/shopspring/decimal"
)

type ISubscribableParser interface {
	Subscribe() (chan *SubTicker, error)
}

// SubTicker represents a ticker that comes from the realtime feeds
type SubTicker struct {
	Ok       bool
	Error    error
	Time     time.Time
	Currency *CurrencyPair
	Value    decimal.Decimal
}

// Subscriber represents a subscribable ticker
type Subscriber struct {
	parser ISubscribableParser
}

// NewSubscriber returns a new subscriber
func NewSubscriber(t SubTickerType) *Subscriber {
	var parser ISubscribableParser
	switch t {
	case CoinoneSubTicker:
		return nil
	case PoloniexSubTicker:
		parser = newPoloniexSubscriber()
	}

	return &Subscriber{parser: parser}
}

func (s *Subscriber) Subscribe() (chan *SubTicker, error) {
	return s.parser.Subscribe()
}
