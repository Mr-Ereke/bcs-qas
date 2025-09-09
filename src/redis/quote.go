package redis

import (
	"context"
	"errors"
	"time"

	"app/models"

	"github.com/redis/go-redis/v9"
	"gitlab.online-fx.com/go-packages/logger"
)

const QuoteChannelPrefix = "QUOTES:"

func (c *Client) ListenChannel(ctx context.Context, symbol string, handler func(quote *models.SymbolQuote)) {
	channel := QuoteChannelPrefix + symbol
	logger.Info("Listening channel: " + channel)

	subscriber := c.client.Subscribe(c.ctx, channel)
	defer func() {
		if err := subscriber.Unsubscribe(ctx, channel); err != nil {
			logger.Errorf("Fail unsubscribe %s channel. Error: %s", channel, err)
		}
		logger.Infof("Unsubscribe %s channel", channel)
		if err := subscriber.Close(); err != nil {
			logger.Errorf("Fail close %s channel subscribe. Error: %s", channel, err)
		}
		logger.Infof("Close %s channel subscribe connection", channel)
	}()

	channelOption := redis.WithChannelHealthCheckInterval(time.Second)
	quoteChannel := subscriber.ChannelWithSubscriptions(channelOption)

	for {
		select {
		case <-ctx.Done():
			return
		case message := <-quoteChannel:
			switch msg := message.(type) {
			case *redis.Subscription:
				logger.Infof("Redis pub/sub %s", msg)
			case *redis.Message:
				symbolQuote, errParse := parseQuote(symbol, msg.Payload)
				if errParse != nil {
					logger.Errorf("Failed parse quote from channel by symbol - %s. Error: %s", symbol, errParse)
					continue
				}
				handler(symbolQuote)
			}
		}
	}
}

func (c *Client) GetQuote(symbol string) (*models.SymbolQuote, error) {
	existsData := c.client.Exists(c.ctx, symbol)
	err := existsData.Err()
	if err != nil {
		return nil, err
	}

	exists := existsData.Val()

	if exists == 0 {
		return nil, ErrNoQuoteData
	}

	data := c.client.Get(c.ctx, symbol)
	err = data.Err()

	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, err
		} else {
			return nil, ErrNoQuoteData
		}
	}

	symbolQuote, parseErr := parseQuote(symbol, data.Val())

	if parseErr != nil {
		logger.Errorf("Failed parse quote from redis by symbol - %s. Error: %s", symbol, parseErr)
	}

	return symbolQuote, nil
}
