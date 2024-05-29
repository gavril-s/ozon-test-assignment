package main

import (
	"log"
	"net/http"
	"ozon-test-assignment/configs"
	"ozon-test-assignment/graph"
	"ozon-test-assignment/internal/comment"
	"ozon-test-assignment/internal/storage"
	"ozon-test-assignment/internal/storage/db"
	"ozon-test-assignment/internal/storage/memory"

	"github.com/99designs/gqlgen/graphql/handler"
)

func main() {
	env := configs.ReadEnv()

	var storageImpl storage.Storage
	if env.StorageType == storage.StorageTypeDB {
		db, err := db.NewDB(env.DatabaseUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		storageImpl = db
	} else if env.StorageType == storage.StorageTypeMemory {
		storageImpl = memory.NewMemory()
	} else {
		log.Fatal("Specified storage type is incorrect")
	}

	commentSubscriptionManager := comment.NewSubscriptionManager()
	resolver := graph.NewResolver(storageImpl, commentSubscriptionManager)
	server := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/query", server)
	log.Fatal(http.ListenAndServe(env.ServerHost+":"+env.ServerPort, nil))
}
