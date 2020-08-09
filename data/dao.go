package data

type DAO interface {
	// Returns true iff subscription with given ID exists.
	Exists(id ID) (bool, error)

	// Returns pointer to COPY of the original data
	GetSubscriptionByID(id ID) (*Subscription, error)

	// Returns pointer to COPY of the original data
	GetAllSubscriptions() (*[]Subscription, error)

	// Returns pointer to COPY of the original data
	GetSubscriptionHistoryByID(id ID) (*History, error)

	// Adds subscription and assigns its own ID.
	// Current ID within subscription is disregarded.
	// Returns pointer to COPY of added subscription
	AddSubscription(subscription *Subscription) (*Subscription, error)

	// Replaces subscription with the pne with matching ID. Returns number of modified documents.
	// That is 1 if subscription was replaced and 0 if matching subscription wasn't found.
	ReplaceSubscription(subscription *Subscription) (int64, error)

	// Deletes subscription with given id.
	// Returns number of deleted objects
	DeleteSubscription(id ID) (int64, error)

	HistoryDAO
}

type HistoryDAO interface {
	// Adds given content to History
	AddToHistory(content *Content) error

	// Deletes all history associated with given Subscription ID
	DeleteAllHistory(subID ID) (int64, error)
}
