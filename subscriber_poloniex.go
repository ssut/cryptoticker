package cryptoticker

import (
	"encoding/json"

	"strings"

	"fmt"

	"time"

	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
)

const (
	poloniexSubURL = "wss://api2.poloniex.com"
)

type poloniexSubscriber struct {
	conn    *websocket.Conn
	markets map[int]*CurrencyPair
	ticker  chan *SubTicker
}

func newPoloniexSubscriber() *poloniexSubscriber {
	ticker := make(chan *SubTicker, 1000)
	subscriber := &poloniexSubscriber{
		ticker:  ticker,
		markets: make(map[int]*CurrencyPair),
	}
	return subscriber
}

func (s *poloniexSubscriber) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(poloniexSubURL, nil)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *poloniexSubscriber) loadMarket() error {
	p := newPoloniexParser()
	t, err := p.RawTicker()
	if err != nil {
		return err
	}

	poloTicker, ok := t.(*poloniexTicker)
	if !ok {
		return fmt.Errorf("could not parse poloniex ticker")
	}

	for name, market := range *poloTicker {
		sep := strings.Split(name, "_")
		pair := &CurrencyPair{sep[0], sep[1]}
		s.markets[market.ID] = pair
	}
	return nil
}

func (s *poloniexSubscriber) Subscribe() (chan *SubTicker, error) {
	if s.conn == nil {
		if err := s.Connect(); err != nil {
			return nil, err
		}
		if err := s.loadMarket(); err != nil {
			return nil, err
		}
	}

	// 1002 = ticker
	s.conn.WriteJSON(map[string]interface{}{"command": "subscribe", "channel": 1002})

	go func() {
		for {
			t := &SubTicker{
				Ok:    false,
				Error: nil,
				Time:  time.Now(),
			}

			_, message, err := s.conn.ReadMessage()
			if err != nil {
				t.Error = err
				s.ticker <- t
				continue
			}

			var data []interface{}
			if err := json.Unmarshal(message, &data); err != nil {
				t.Error = err
				s.ticker <- t
				continue
			}
			// to filter out the heartbeat traffics
			if len(data) < 3 {
				continue
			}

			items := data[2].([]interface{})
			marketID := int(items[0].(float64))
			var pair *CurrencyPair
			var ok bool
			if pair, ok = s.markets[marketID]; !ok {
				t.Error = fmt.Errorf("failed to find an currency")
				goto sendErrorAndContinue
			}

			t.Currency = pair
			t.Value, err = decimal.NewFromString(items[1].(string))
			if err != nil {
				goto sendErrorAndContinue
			}

			t.Ok = true
			s.ticker <- t
			continue

		sendErrorAndContinue:
			s.ticker <- t
			continue
		}
	}()

	return s.ticker, nil
}
