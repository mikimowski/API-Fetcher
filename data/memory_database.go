package data

import (
	"errors"
	"sync"
)

type data struct {
	subscriptions map[ID]*Subscription
	history       map[ID]History
}

// Memory Database that implements DAO interface
// Provides access to underlying data
type MemoryDB struct {
	data   data
	mtx    sync.RWMutex
	nextID ID
}

func NewMemoryDB() (*MemoryDB, error) {
	return &MemoryDB{
		data:   data{
			subscriptions: map[ID]*Subscription{},
			history: map[ID]History{},
		},
		mtx:    sync.RWMutex{},
		nextID: 1,
	}, nil
}

// Not thread safe. Caller is responsible for acquiring mutex!
func (db *MemoryDB) reserveID() (id ID) {
	id = db.nextID
	db.nextID++
	return
}

//************* reading data *************//

func (db *MemoryDB) Exists(id ID) (bool, error) {
	db.mtx.RLock()
	defer db.mtx.RUnlock()

	_, ok := db.data.subscriptions[id]
	return ok, nil
}

func (db *MemoryDB) GetSubscriptionByID(id ID) (*Subscription, error) {
	db.mtx.RLock()
	defer db.mtx.RUnlock()

	if sub, ok := db.data.subscriptions[id]; ok {
		subCopy := *sub
		return &subCopy, nil
	}
	return nil, errors.New("subscription not found")
}

func (db *MemoryDB) GetAllSubscriptions() (*[]Subscription, error) {
	db.mtx.RLock()
	defer db.mtx.RUnlock()

	subs := &[]Subscription{}
	for _, v := range db.data.subscriptions {
		*subs = append(*subs, *v)
	}
	return subs, nil
}

func (db *MemoryDB) GetSubscriptionHistoryByID(id ID) (*History, error) {
	db.mtx.RLock()
	defer db.mtx.RUnlock()

	h := &History{}
	for _, c := range db.data.history[id] {
		*h = append(*h, *c.DeepCopy())
	}
	return h, nil
}

//************* altering data *************//

func (db *MemoryDB) AddSubscription(subscription *Subscription) (*Subscription, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	id := db.reserveID()
	sub := Subscription{
		ID:       id,
		URL:      subscription.URL,
		Interval: subscription.Interval,
	}
	db.data.subscriptions[id] = &sub

	subCopy := sub
	return &subCopy, nil
}

func (db *MemoryDB) ReplaceSubscription(subscription *Subscription) (int64, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	if sub, ok := db.data.subscriptions[subscription.ID]; ok {
		subCopy := *sub
		db.data.subscriptions[subscription.ID] = &subCopy
		return 1, nil
	}
	return 0, nil
}

func (db *MemoryDB) DeleteSubscription(id ID) (int64, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	if _, ok := db.data.subscriptions[id]; !ok {
		return 0, nil
	}
	delete(db.data.subscriptions, id)
	return 1, nil
}

//************* HistoryDAO interface *************//

func (db *MemoryDB) AddToHistory(content *Content) error {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	db.data.history[content.SubID] = append(db.data.history[content.SubID], *content.DeepCopy())
	return nil
}

func (db *MemoryDB) DeleteAllHistory(subID ID) (int64, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	if _, ok := db.data.history[subID]; !ok {
		return 0, nil
	}
	cnt := len(db.data.history[subID])
	delete(db.data.history, subID)
	return int64(cnt), nil
}
