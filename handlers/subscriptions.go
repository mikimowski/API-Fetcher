package handlers

import (
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/data"
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/subscriber"
	"go.uber.org/zap"
)

type Subscriptions struct {
	dao        data.DAO
	subscriber *subscriber.Subscriber
	l          *zap.SugaredLogger
}

func NewSubscriptions(dao data.DAO, l *zap.SugaredLogger, s *subscriber.Subscriber) *Subscriptions {
	return &Subscriptions{
		dao:        dao,
		subscriber: s,
		l:          l,
	}
}
