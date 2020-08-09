package handlers

import (
	"github.com/go-chi/chi"
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/mock"
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/subscriber"
)

// Basic handler based on mock memory database.
func getHandlerMockMemoryDB() *Subscriptions {
	return &Subscriptions{
		dao:        mock.NewMemoryDatabase(),
		subscriber: subscriber.NewSubscriber(mock.NewMemoryDatabase(), mock.Logger.Sugar()),
		l:          mock.Logger.Sugar(),
	}
}

func getHandlerFailingDB() *Subscriptions {
	return &Subscriptions{
		dao:        &mock.FailingDB{},
		subscriber: subscriber.NewSubscriber(mock.NewMemoryDatabase(), mock.Logger.Sugar()),
		l:          mock.Logger.Sugar(),
	}
}

func getApiFetcherChiRouter(sh *Subscriptions) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/api/fetcher", func(r chi.Router) {
		r.Use(sh.ContentTypeJSON)
		r.With(sh.PayloadLimit).Post("/", sh.Add)
		r.Get("/", sh.ListAll)

		r.Route("/{id:[0-9]+}", func(r chi.Router) {
			r.Use(sh.IDContext)
			r.Get("/history", sh.ListHistory)
			r.Delete("/", sh.Delete)
			r.With(sh.PayloadLimit).Patch("/", sh.Update)
		})
	})
	return r
}
