package models

import (
	"fmt"
	"sync"
)

type PriceTypeDirection struct {
	mutex      sync.RWMutex
	BidUp      map[uint]*Subscription
	BidDown    map[uint]*Subscription
	AskUp      map[uint]*Subscription
	AskDown    map[uint]*Subscription
	LastUp     map[uint]*Subscription
	LastDown   map[uint]*Subscription
	TypeExists *TypeExists
}

const (
	Bid  = "Bid"
	Ask  = "Ask"
	Last = "Last"
)

func NewPriceTypeDirection() *PriceTypeDirection {
	return &PriceTypeDirection{
		BidUp:      make(map[uint]*Subscription),
		BidDown:    make(map[uint]*Subscription),
		AskUp:      make(map[uint]*Subscription),
		AskDown:    make(map[uint]*Subscription),
		LastUp:     make(map[uint]*Subscription),
		LastDown:   make(map[uint]*Subscription),
		TypeExists: NewTypeExists(),
	}
}

func (ptd *PriceTypeDirection) Get(priceType string, direction bool) (map[uint]*Subscription, error) {
	ptd.mutex.RLock()
	defer ptd.mutex.RUnlock()

	switch priceType {
	case Bid:
		if direction {
			return ptd.BidUp, nil
		} else {
			return ptd.BidDown, nil
		}
	case Ask:
		if direction {
			return ptd.AskUp, nil
		} else {
			return ptd.AskDown, nil
		}
	case Last:
		if direction {
			return ptd.LastUp, nil
		} else {
			return ptd.LastDown, nil
		}
	}

	return map[uint]*Subscription{}, fmt.Errorf("unknow price type for get. Type: %s", priceType)
}

func (ptd *PriceTypeDirection) Create(subscription *Subscription) error {
	ptd.mutex.Lock()
	defer ptd.mutex.Unlock()

	switch subscription.PriceType {
	case Bid:
		if subscription.PriceDirection {
			if _, exists := ptd.BidUp[subscription.ID]; exists {
				return fmt.Errorf("already exists alert %s Up", Bid)
			}
			ptd.BidUp[subscription.ID] = subscription
			if !ptd.TypeExists.IsExists(Bid, true) {
				ptd.TypeExists.SetExists(true, Bid, true)
			}
		} else {
			if _, exists := ptd.BidUp[subscription.ID]; exists {
				return fmt.Errorf("already exists alert %s Down", Bid)
			}
			ptd.BidDown[subscription.ID] = subscription
			if !ptd.TypeExists.IsExists(Bid, false) {
				ptd.TypeExists.SetExists(true, Bid, false)
			}
		}
	case Ask:
		if subscription.PriceDirection {
			if _, exists := ptd.AskUp[subscription.ID]; exists {
				return fmt.Errorf("already exists alert %s Up", Ask)
			}
			ptd.AskUp[subscription.ID] = subscription
			if !ptd.TypeExists.IsExists(Ask, true) {
				ptd.TypeExists.SetExists(true, Ask, true)
			}
		} else {
			if _, exists := ptd.AskDown[subscription.ID]; exists {
				return fmt.Errorf("already exists alert %s Down", Ask)
			}
			ptd.AskDown[subscription.ID] = subscription
			if !ptd.TypeExists.IsExists(Ask, false) {
				ptd.TypeExists.SetExists(true, Ask, false)
			}
		}
	case Last:
		if subscription.PriceDirection {
			if _, exists := ptd.LastUp[subscription.ID]; exists {
				return fmt.Errorf("already exists alert %s Up", Last)
			}
			ptd.LastUp[subscription.ID] = subscription
			if !ptd.TypeExists.IsExists(Last, true) {
				ptd.TypeExists.SetExists(true, Last, true)
			}
		} else {
			if _, exists := ptd.LastDown[subscription.ID]; exists {
				return fmt.Errorf("already exists alert %s Down", Last)
			}
			ptd.LastDown[subscription.ID] = subscription
			if !ptd.TypeExists.IsExists(Last, false) {
				ptd.TypeExists.SetExists(true, Last, false)
			}
		}
	default:
		return fmt.Errorf("unknow price type for create. Type: %s", subscription.PriceType)
	}
	return nil
}

func (ptd *PriceTypeDirection) Delete(subscription *Subscription) error {
	ptd.mutex.Lock()
	defer ptd.mutex.Unlock()

	switch subscription.PriceType {
	case Bid:
		if subscription.PriceDirection {
			if _, exists := ptd.BidUp[subscription.ID]; !exists {
				return fmt.Errorf("not found alert %s Up", Bid)
			}
			delete(ptd.BidUp, subscription.ID)
			if len(ptd.BidUp) == 0 {
				ptd.TypeExists.SetExists(false, Bid, true)
			}
		} else {
			if _, exists := ptd.BidDown[subscription.ID]; !exists {
				return fmt.Errorf("not found alert %s Down", Bid)
			}
			delete(ptd.BidDown, subscription.ID)
			if len(ptd.BidDown) == 0 {
				ptd.TypeExists.SetExists(false, Bid, false)
			}
		}
	case Ask:
		if subscription.PriceDirection {
			if _, exists := ptd.AskUp[subscription.ID]; !exists {
				return fmt.Errorf("not found alert %s Up", Ask)
			}
			delete(ptd.AskUp, subscription.ID)
			if len(ptd.AskUp) == 0 {
				ptd.TypeExists.SetExists(false, Ask, true)
			}
		} else {
			if _, exists := ptd.AskDown[subscription.ID]; !exists {
				return fmt.Errorf("not found alert %s Down", Ask)
			}
			delete(ptd.AskDown, subscription.ID)
			if len(ptd.AskDown) == 0 {
				ptd.TypeExists.SetExists(false, Ask, false)
			}
		}
	case Last:
		if subscription.PriceDirection {
			if _, exists := ptd.LastUp[subscription.ID]; !exists {
				return fmt.Errorf("not found alert %s Up", Last)
			}
			delete(ptd.LastUp, subscription.ID)
			if len(ptd.LastUp) == 0 {
				ptd.TypeExists.SetExists(false, Last, true)
			}
		} else {
			if _, exists := ptd.LastDown[subscription.ID]; !exists {
				return fmt.Errorf("not found alert %s Down", Last)
			}
			delete(ptd.LastDown, subscription.ID)
			if len(ptd.LastDown) == 0 {
				ptd.TypeExists.SetExists(false, Last, false)
			}
		}
	default:
		return fmt.Errorf("unknow price type for delete. Type: %s", subscription.PriceType)
	}

	return nil
}

func (ptd *PriceTypeDirection) GetSubscriptionListByCustomerId(customerId uint) []*Subscription {
	ptd.mutex.RLock()
	defer ptd.mutex.RUnlock()

	customerList := make([]*Subscription, 0, ptd.GetCount())

	for _, alert := range ptd.getAll() {
		if customerId == alert.CustomerId {
			customerList = append(customerList, alert)
		}
	}

	return customerList
}

func (ptd *PriceTypeDirection) GetSubscriptionCountByCustomerId(customerId uint) int {
	ptd.mutex.RLock()
	defer ptd.mutex.RUnlock()

	count := 0

	for _, alert := range ptd.getAll() {
		if customerId == alert.CustomerId {
			count++
		}
	}

	return count
}

func (ptd *PriceTypeDirection) getAll() []*Subscription {
	ptd.mutex.RLock()
	defer ptd.mutex.RUnlock()

	list := make([]*Subscription, 0, ptd.GetCount())

	for _, alert := range ptd.BidUp {
		list = append(list, alert)
	}
	for _, alert := range ptd.BidDown {
		list = append(list, alert)
	}
	for _, alert := range ptd.AskUp {
		list = append(list, alert)
	}
	for _, alert := range ptd.AskDown {
		list = append(list, alert)
	}
	for _, alert := range ptd.LastUp {
		list = append(list, alert)
	}
	for _, alert := range ptd.LastDown {
		list = append(list, alert)
	}

	return list
}

func (ptd *PriceTypeDirection) GetCount() int {
	ptd.mutex.RLock()
	defer ptd.mutex.RUnlock()

	count := 0

	count += len(ptd.BidUp)
	count += len(ptd.BidDown)
	count += len(ptd.AskUp)
	count += len(ptd.AskDown)
	count += len(ptd.LastUp)
	count += len(ptd.LastDown)

	return count
}
