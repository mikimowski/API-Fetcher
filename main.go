package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/data"
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/handlers"
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/subscriber"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"time"
)

const serverAddr = ":8080"

// To work locally set ENV: "mongodb://localhost:27017"
const mongoEnvKey = "MONGO_URL"
const apiFetcherDatabase = "api-fetcher"

// Creates routing for api-fetcher
func newApiFetcherRouter(dao data.DAO, logger *zap.SugaredLogger) *chi.Mux {
	sh := handlers.NewSubscriptions(dao, logger, subscriber.NewSubscriber(dao, logger))

	r := chi.NewRouter()
	// Router for /api/fetcher resources
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

func newLogger() (*zap.Logger, error) {
	if len(os.Args) > 1 && os.Args[1] == "debug" {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}

func initMongoConnection(mongoConnectionURL string) (data.DAO, *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Connect with mongo
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConnectionURL))
	if err != nil {
		log.Fatalf("error connecting to mongodb: %s", err.Error())
	}
	// Check if succeeded
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("error connecting to mongodb: %s", err.Error())
	}

	// Development hack for convenience ;). Allows to reuse running container
	db := client.Database(apiFetcherDatabase)
	err = db.Drop(ctx)
	if err != nil {
		log.Fatalf("error initiating mongo: %s", err.Error())
	}
	db = client.Database(apiFetcherDatabase)

	// Create MongoDAO for created database
	mongoDAO := data.NewMongoDAO(db, context.Background())
	if err := mongoDAO.Init(); err != nil {
		log.Fatalf("error initiating MongoDAO: %s", err.Error())
	}

	return mongoDAO, client
}

func initServer(dao data.DAO) *http.Server {
	logger, err := newLogger()
	if err != nil {
		log.Fatalf("error initiating logger: %s", err)
	}

	r := newApiFetcherRouter(dao, logger.Sugar())

	return &http.Server{
		Addr:         serverAddr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

func main() {
	var dao data.DAO
	if mongoConnectionURL := os.Getenv(mongoEnvKey); mongoConnectionURL != "" {
		var client *mongo.Client
		dao, client = initMongoConnection(mongoConnectionURL)
		defer func() {
			if err := client.Disconnect(context.Background()); err != nil {
				panic(err)
			}
		}()
	} else {
		dao, _ = data.NewMemoryDB()
	}

	server := initServer(dao)
	log.Fatal(server.ListenAndServe())
}
