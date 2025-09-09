package models

import "sync"

type TypeExists struct {
	mutex          sync.RWMutex
	BidUpExists    bool
	BidDownExists  bool
	AskUpExists    bool
	AskDownExists  bool
	LastUpExists   bool
	LastDownExists bool
}

func NewTypeExists() *TypeExists {
	return &TypeExists{
		BidUpExists:    false,
		BidDownExists:  false,
		AskUpExists:    false,
		AskDownExists:  false,
		LastUpExists:   false,
		LastDownExists: false,
	}
}

func (te *TypeExists) SetExists(exists bool, priceType string, direction bool) {
	te.mutex.Lock()
	defer te.mutex.Unlock()

	switch priceType {
	case Bid:
		if direction {
			te.BidUpExists = exists
		} else {
			te.BidDownExists = exists
		}
	case Ask:
		if direction {
			te.AskUpExists = exists
		} else {
			te.AskDownExists = exists
		}
	case Last:
		if direction {
			te.LastUpExists = exists
		} else {
			te.LastDownExists = exists
		}
	}
}

func (te *TypeExists) IsExists(priceType string, direction bool) bool {
	te.mutex.RLock()
	defer te.mutex.RUnlock()

	switch priceType {
	case Bid:
		if direction {
			return te.BidUpExists
		} else {
			return te.BidDownExists
		}
	case Ask:
		if direction {
			return te.AskUpExists
		} else {
			return te.AskDownExists
		}
	case Last:
		if direction {
			return te.LastUpExists
		} else {
			return te.LastDownExists
		}
	}

	return false
}

func (te *TypeExists) IsEmpty() bool {
	te.mutex.RLock()
	defer te.mutex.RUnlock()

	if te.BidUpExists {
		return false
	}

	if te.BidDownExists {
		return false
	}

	if te.AskUpExists {
		return false
	}

	if te.AskDownExists {
		return false
	}

	if te.LastUpExists {
		return false
	}

	if te.LastDownExists {
		return false
	}

	return true
}
