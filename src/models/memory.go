package models

import "sync"

type Memory struct {
	mutex  sync.RWMutex
	Alerts map[string]*PriceTypeDirection
}

func NewMemory() *Memory {
	return &Memory{Alerts: make(map[string]*PriceTypeDirection)}
}

func (m *Memory) Get(instrument string) (*PriceTypeDirection, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	entity, exist := m.Alerts[instrument]

	return entity, exist
}

func (m *Memory) Set(instrument string, entity *PriceTypeDirection) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.Alerts[instrument] = entity
}

func (m *Memory) Delete(instrument string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.Alerts, instrument)
}
