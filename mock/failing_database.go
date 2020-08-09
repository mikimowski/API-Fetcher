// Always failing database

package mock

import (
	"errors"
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/data"
)

type FailingDB struct{}

//************* reading data *************//

func (db *FailingDB) Exists(id data.ID) (bool, error) {
	return false, errors.New("dummy error")
}

func (db *FailingDB) GetSubscriptionByID(id data.ID) (*data.Subscription, error) {
	return nil, errors.New("subscription not found")
}

func (db *FailingDB) GetAllSubscriptions() (*[]data.Subscription, error) {
	return nil, errors.New("dummy error")
}

func (db *FailingDB) GetSubscriptionHistoryByID(id data.ID) (*data.History, error) {
	return nil, errors.New("dummy error")
}

//************* altering data *************//

func (db *FailingDB) AddSubscription(subscription *data.Subscription) (*data.Subscription, error) {
	return nil, errors.New("dummy error")
}

func (db *FailingDB) ReplaceSubscription(subscription *data.Subscription) (int64, error) {
	return 0, errors.New("dummy error")
}

func (db *FailingDB) DeleteSubscription(id data.ID) (int64, error) {
	return 0, errors.New("dummy error")
}

// HistoryDAO interface
func (db *FailingDB) AddToHistory(content *data.Content) error {
	return errors.New("dummy error")
}

func (db *FailingDB) DeleteAllHistory(subID data.ID) (int64, error) {
	return 0, errors.New("dummy error")
}
