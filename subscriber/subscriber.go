// Subscriber is used to organize fetching subscribed urls.
// Each url is handled by at most one goroutine.

// Subscriber stores mapping: subscriptionID -> chan command
// those channels are distributed to workers (goroutines) responsible for given subscription.
// Each channel is given to exactly one worker.
// On update channel for given subscription is overwritten.

// Subscriber is thread safe
// Unsubscribe:

package subscriber

import (
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/data"
	"go.uber.org/zap"
	"sync"
)

type command int

const stop = command(0)
const stopAndClean = command(1)

type Subscriber struct {
	// Those channels are used to communicate with goroutines responsible for given subscription
	stopChan map[data.ID]chan command
	dao      data.HistoryDAO
	l        *zap.SugaredLogger
	mtx      sync.Mutex
}

func NewSubscriber(dao data.DAO, l *zap.SugaredLogger) *Subscriber {
	return &Subscriber{
		stopChan: make(map[data.ID]chan command),
		dao:      dao,
		l:        l,
		mtx:      sync.Mutex{},
	}
}

// First call for subscription with specific ID setups worker and runs it.
// Subsequent calls with subscription with given ID won't have any effect.
// To make changes in url or interval use Update.
func (s *Subscriber) Subscribe(sub *data.Subscription) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if _, ok := s.stopChan[sub.ID]; ok {
		return
	}
	// add subscription
	// BUFFERED channel, which is crucial for performance ~ notify and DO NOT wait.
	stopChan := make(chan command, 1)
	s.stopChan[sub.ID] = stopChan

	w := worker{
		sub:      *sub,
		dao:      s.dao,
		l:        s.l,
		stopChan: stopChan,
	}
	go w.run()
}

// Use to update url or interval for given subscription.
// Assumption: subscription with given ID should have been register via Subscribe before using update.
// If subscription with given ID hasn't been registered then this call will have no effect.
// Under the hood:
// Stops current worker and subscribes with new settings
// It's similar to calling Unsubscribe and then Subscribe but here worker will not delete history.
// Moreover, mutex is acquired once.
func (s *Subscriber) Update(sub *data.Subscription) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if _, ok := s.stopChan[sub.ID]; !ok {
		return
	}
	// Unsubscribe
	s.stopChan[sub.ID] <- stop
	delete(s.stopChan, sub.ID)
	// Resubscribe with new parameters
	stopChan := make(chan command, 1)
	s.stopChan[sub.ID] = stopChan

	w := worker{
		sub:      *sub,
		dao:      s.dao,
		l:        s.l,
		stopChan: stopChan,
	}
	go w.run()
}

// Subscriber will stop following URL immediately.
// History will be removed but NOT immediately. Adequate worker will delete history associated with given subscription.
// Subscriber only notifies him about this.
func (s *Subscriber) Unsubscribe(subID data.ID) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if stopChan, ok := s.stopChan[subID]; ok {
		stopChan <- stopAndClean
		delete(s.stopChan, subID)
	}
}
