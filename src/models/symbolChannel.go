package models

import (
	"context"
	"sync"
)

type SymbolChannel struct {
	mutex    sync.RWMutex
	channels map[string]context.CancelFunc
}

func NewSymbolChannel() *SymbolChannel {
	return &SymbolChannel{channels: make(map[string]context.CancelFunc)}
}

func (sc *SymbolChannel) Set(symbol string, cancel context.CancelFunc) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	sc.channels[symbol] = cancel
}

func (sc *SymbolChannel) Get(symbol string) (context.CancelFunc, bool) {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	cancel, exists := sc.channels[symbol]
	return cancel, exists
}
