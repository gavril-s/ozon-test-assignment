package main

import (
	"log"
	"net/http"
	"os"
	"ozon-test-assignment/graph"
	"ozon-test-assignment/internal/comment"
	"ozon-test-assignment/internal/storage"
	"ozon-test-assignment/internal/storage/db"
	"ozon-test-assignment/internal/storage/memory"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultHost = "0.0.0.0"
const defaultPort = "8080"
const defaultStorageType = storage.StorageTypeDB
const postgresUser = "user"
const postgresPassword = "password"
const postgresHost = "localhost"
const postgresPort = "5432"
const postgresDB = "db"

type envData struct {
	host        string
	port        string
	storageType string
}

func readEnv() envData {
	host := os.Getenv("PORT")
	if host == "" {
		host = defaultHost
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = defaultStorageType
	}

	return envData{
		host:        host,
		port:        port,
		storageType: storageType,
	}
}

func main() {
	env := readEnv()

	var storageImpl storage.Storage
	if env.storageType == storage.StorageTypeDB {
		db, err := db.NewDB(postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB)
		if err != nil {
			log.Fatal("Database connection failed")
		}
		defer db.Close()
		storageImpl = db
	} else if env.storageType == storage.StorageTypeMemory {
		storageImpl = memory.NewMemory()
	} else {
		log.Fatal("Specified storage type is incorrect")
	}

	commentSubscriptionManager := comment.NewSubscriptionManager()
	resolver := graph.NewResolver(storageImpl, commentSubscriptionManager)
	server := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", server)

	log.Printf("connect to http://%s:%s/ for GraphQL playground", env.host, env.port)
	log.Fatal(http.ListenAndServe(env.host+":"+env.port, nil))
}
