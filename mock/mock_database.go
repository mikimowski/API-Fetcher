package mock

import (
	"errors"
	data "github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/data"
	"sync"
)

type mockData struct {
	subscriptions map[data.ID]*data.Subscription
	history       map[data.ID]data.History
}

// Memory Database that implements data.DAO interface
// Provides access to underlying data
type MemoryDatabase struct {
	data   mockData
	mtx    sync.RWMutex
	nextID data.ID
}

// Few records initially stored in DB for testing purposes.

var s1 = "my mock history"

var cont1 = data.Content{
	SubID:     1,
	Response:  &s1,
	Duration:  0.532,
	CreatedAt: "1559034938.638",
}

var cont2 = data.Content{
	SubID:     1,
	Response:  nil,
	Duration:  5,
	CreatedAt: "1559034938.638",
}

var sub1 = data.Subscription{
	ID:       1,
	URL:      "https://httpbin.org/range/15",
	Interval: 60,
}

var sub2 = data.Subscription{
	ID:       2,
	URL:      "https://httpbin.org/delay/10",
	Interval: 120,
}

var history1 = data.History{
	cont1,
	cont2,
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		data: mockData{
			subscriptions: map[data.ID]*data.Subscription{
				1: &sub1,
				2: &sub2,
			},
			history: map[data.ID]data.History{
				1: history1,
			},
		},
		mtx:    sync.RWMutex{},
		nextID: 3,
	}
}

// Not thread safe. Caller is responsible for acquiring mutex!
func (db *MemoryDatabase) reserveID() (id data.ID) {
	id = db.nextID
	db.nextID++
	return
}

//************* reading data *************//

func (db *MemoryDatabase) Exists(id data.ID) (bool, error) {
	db.mtx.RLock()
	defer db.mtx.RUnlock()

	_, ok := db.data.subscriptions[id]
	return ok, nil
}

func (db *MemoryDatabase) GetSubscriptionByID(id data.ID) (*data.Subscription, error) {
	db.mtx.RLock()
	defer db.mtx.RUnlock()

	if sub, ok := db.data.subscriptions[id]; ok {
		subCopy := *sub
		return &subCopy, nil
	}
	return nil, errors.New("subscription not found")
}

func (db *MemoryDatabase) GetAllSubscriptions() (*[]data.Subscription, error) {
	db.mtx.RLock()
	defer db.mtx.RUnlock()

	subs := &[]data.Subscription{}
	for _, v := range db.data.subscriptions {
		*subs = append(*subs, *v)
	}
	return subs, nil
}

func (db *MemoryDatabase) GetSubscriptionHistoryByID(id data.ID) (*data.History, error) {
	db.mtx.RLock()
	defer db.mtx.RUnlock()

	h := &data.History{}
	for _, c := range db.data.history[id] {
		*h = append(*h, *c.DeepCopy())
	}
	return h, nil
}

//************* altering data *************//

func (db *MemoryDatabase) AddSubscription(subscription *data.Subscription) (*data.Subscription, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	id := db.reserveID()
	sub := data.Subscription{
		ID:       id,
		URL:      subscription.URL,
		Interval: subscription.Interval,
	}
	db.data.subscriptions[id] = &sub

	subCopy := sub
	return &subCopy, nil
}

func (db *MemoryDatabase) ReplaceSubscription(subscription *data.Subscription) (int64, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	if sub, ok := db.data.subscriptions[subscription.ID]; ok {
		subCopy := *sub
		db.data.subscriptions[subscription.ID] = &subCopy
		return 1, nil
	}
	return 0, nil
}

func (db *MemoryDatabase) DeleteSubscription(id data.ID) (int64, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	if _, ok := db.data.subscriptions[id]; !ok {
		return 0, nil
	}
	delete(db.data.subscriptions, id)
	return 1, nil
}

// HistoryDAO interface
func (db *MemoryDatabase) AddToHistory(content *data.Content) error {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	db.data.history[content.SubID] = append(db.data.history[content.SubID], *content.DeepCopy())
	return nil
}

func (db *MemoryDatabase) DeleteAllHistory(subID data.ID) (int64, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	if _, ok := db.data.history[subID]; !ok {
		return 0, nil
	}
	cnt := len(db.data.history[subID])
	delete(db.data.history, subID)
	return int64(cnt), nil
}
