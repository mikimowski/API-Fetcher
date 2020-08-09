package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const subIDSequenceKey = "subscriptionID"
const countersCollection = "counters"
const subscriptionsCollection = "subscriptions"
const historyCollection = "history"
const startID = 1 // ID == 0 can be thought of as uninitialized

type MongoDAO struct {
	client *mongo.Client
	db     *mongo.Database
	ctx    context.Context
}

func NewMongoDAO(db *mongo.Database, ctx context.Context) *MongoDAO {
	return &MongoDAO{
		db:  db,
		ctx: ctx,
	}
}

//************* miscellaneous *************//

// Initialize custom autoIncrement
func (dao *MongoDAO) Init() error {
	_, err := dao.db.Collection(countersCollection).InsertOne(dao.ctx, bson.M{"_id": subIDSequenceKey, "value": ID(startID)})
	if err != nil {
		return err
	}
	return nil
}

// Custom autoIncrement. Solution suggested by mongodb documentation.
func (dao *MongoDAO) getNextID() (ID, error) {
	countersCollection := dao.db.Collection(countersCollection)

	// Returns document from before update
	filter := bson.M{"_id": subIDSequenceKey}
	update := bson.D{{"$inc", bson.M{"value": 1}}}
	singleResult := countersCollection.FindOneAndUpdate(dao.ctx, filter, update)
	var result struct {
		Value int64
	}
	if err := singleResult.Decode(&result); err != nil {
		return 0, err
	}
	return ID(result.Value), nil
}

//************* reading data *************//

func (dao *MongoDAO) Exists(id ID) (bool, error) {
	subscriptionsCollection := dao.db.Collection(subscriptionsCollection)

	var sub Subscription
	filter := bson.M{"_id": id}
	if err := subscriptionsCollection.FindOne(dao.ctx, filter).Decode(&sub); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (dao *MongoDAO) GetSubscriptionByID(id ID) (*Subscription, error) {
	subscriptionsCollection := dao.db.Collection(subscriptionsCollection)

	var sub Subscription
	filter := bson.M{"_id": id}
	if err := subscriptionsCollection.FindOne(dao.ctx, filter).Decode(&sub); err != nil {
		return nil, err
	}
	return &sub, nil
}

func (dao *MongoDAO) GetAllSubscriptions() (*[]Subscription, error) {
	subscriptionsCollection := dao.db.Collection(subscriptionsCollection)

	filter := bson.M{}
	cursor, err := subscriptionsCollection.Find(dao.ctx, filter)
	if err != nil {
		return nil, err
	}

	subs := &[]Subscription{}
	if err = cursor.All(dao.ctx, subs); err != nil {
		return nil, err
	}

	return subs, nil
}

func (dao *MongoDAO) GetSubscriptionHistoryByID(id ID) (*History, error) {
	historyCollection := dao.db.Collection(historyCollection)

	filter := bson.M{"subID": id}
	cursor, err := historyCollection.Find(dao.ctx, filter)
	if err != nil {
		return nil, err
	}

	history := &History{}
	if err = cursor.All(dao.ctx, history); err != nil {
		return nil, err
	}
	return history, nil
}

//************* altering data *************//

func (dao *MongoDAO) AddSubscription(subscription *Subscription) (*Subscription, error) {
	subscriptionsCollection := dao.db.Collection(subscriptionsCollection)

	id, err := dao.getNextID()
	if err != nil {
		return nil, err
	}
	sub := Subscription{
		ID:       id,
		URL:      subscription.URL,
		Interval: subscription.Interval,
	}
	if _, err := subscriptionsCollection.InsertOne(dao.ctx, sub); err != nil {
		return nil, err
	}
	return &sub, nil
}

func (dao *MongoDAO) ReplaceSubscription(subscription *Subscription) (int64, error) {
	subscriptionsCollection := dao.db.Collection(subscriptionsCollection)

	filter := bson.M{"_id": subscription.ID}
	updateResult, err := subscriptionsCollection.ReplaceOne(dao.ctx, filter, subscription)
	if err != nil {
		return 0, err
	}
	return updateResult.ModifiedCount, nil
}

func (dao *MongoDAO) DeleteSubscription(id ID) (int64, error) {
	subscriptionsCollection := dao.db.Collection(subscriptionsCollection)

	filter := bson.M{"_id": id}
	deleteResult, err := subscriptionsCollection.DeleteOne(dao.ctx, filter)
	if err != nil {
		return 0, err
	}
	return deleteResult.DeletedCount, nil
}

//************* HistoryDAO interface *************//

func (dao *MongoDAO) AddToHistory(content *Content) error {
	historyCollection := dao.db.Collection(historyCollection)

	_, err := historyCollection.InsertOne(dao.ctx, content)
	if err != nil {
		return err
	}
	return nil
}

func (dao *MongoDAO) DeleteAllHistory(subID ID) (int64, error) {
	historyCollection := dao.db.Collection(historyCollection)

	filter := bson.M{"subID": subID}
	deleteResult, err := historyCollection.DeleteMany(dao.ctx, filter)

	if err != nil {
		return 0, err
	}
	return deleteResult.DeletedCount, nil
}
