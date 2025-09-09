package service

import (
	"context"

	"app/models"
	"app/redis"

	"gitlab.online-fx.com/go-packages/gormdb"
	"gitlab.online-fx.com/go-packages/logger"
	"gorm.io/gorm"
)

var (
	ctx           context.Context  //контекст
	settings      *models.Settings //для валидации
	redisClient   *redis.Client
	symbolChannel *models.SymbolChannel
	memory        *models.Memory
	db            *gorm.DB
)

func Load(config *models.Config, redis *redis.Client) {
	ctx = context.Background()
	settings = &config.Settings
	redisClient = redis
	symbolChannel = models.NewSymbolChannel()
	db = gormdb.GetClient(models.ServiceDB)
	memory = models.NewMemory()
}

func Restore() {
	var subscriptionList []*models.Subscription

	db = gormdb.GetClient(models.ServiceDB)
	err := db.Model(&models.Subscription{}).Find(&subscriptionList).Error

	if err != nil {
		logger.Fatalf("Failed get alerts. DB error: %s", err)
	}

	for _, subscription := range subscriptionList {
		err = createSubscription(subscription)
		if err != nil {
			logger.Fatalf("Failed load alerts. Service error: %s", err)
		}
	}
}

func SubscribeInstrument(symbol string) {
	cancelCtx, cancel := context.WithCancel(ctx)
	symbolChannel.Set(symbol, cancel)

	go redisClient.ListenChannel(cancelCtx, symbol, handler)
}

func UnsubscribeInstrument(symbol string) {
	cancelFunc, exist := symbolChannel.Get(symbol)

	if !exist {
		logger.Errorf("Can not unsubscribe %s symbol channel", symbol)
		return
	}

	cancelFunc()
}

func handler(quote *models.SymbolQuote) {
	priceTypeDirection, instrumentExists := memory.Get(quote.Symbol)

	if !instrumentExists {
		logger.Errorf("Not found instrument in memory for handle. Symbol: %s", quote.Symbol)
		return
	}

	if priceTypeDirection.TypeExists.IsExists(models.Bid, true) {
		// Обработка Bid Up
		go upHandler(quote, models.Bid)
	}

	if priceTypeDirection.TypeExists.IsExists(models.Bid, false) {
		// Обработка Bid Down
		go downHandler(quote, models.Bid)
	}

	if priceTypeDirection.TypeExists.IsExists(models.Ask, true) {
		// Обработка Ask Up
		go upHandler(quote, models.Ask)
	}

	if priceTypeDirection.TypeExists.IsExists(models.Ask, false) {
		// Обработка Ask Down
		go downHandler(quote, models.Ask)
	}

	if priceTypeDirection.TypeExists.IsExists(models.Last, true) {
		// Обработка Last Up
		go upHandler(quote, models.Last)
	}

	if priceTypeDirection.TypeExists.IsExists(models.Last, false) {
		// Обработка Last Down
		go downHandler(quote, models.Last)
	}
}
