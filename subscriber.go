package cryptoticker

type ISubscribableParser interface
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
		parser = newCoinoneSubscriber()
	case PoloniexSubTicker:
		parser = newPoloniexSubscriber()
	}

	return &Subscriber{parser: parser}
}
