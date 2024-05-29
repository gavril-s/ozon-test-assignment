package configs

import (
	"os"
	"ozon-test-assignment/internal/storage"
)

const CommentsContentLengthLimit = 2000
const CommentsDepthLimit = 3
const CommentsThreadLimit = 100

const PostSnippetDefaultLength = 300
const PostSnippetsLimit = 100

const DefaultServerHost = "0.0.0.0"
const DefaultServerPort = "8080"
const DefaultStorageType = storage.StorageTypeMemory
const DefaultDatabaseUrl = "user=user password=password host=postgres port=5432 " +
	"dbname=db sslmode=disable TimeZone=Europe/Moscow"

type Env struct {
	ServerHost  string
	ServerPort  string
	StorageType string
	DatabaseUrl string
}

func ReadEnv() Env {
	host := os.Getenv("HOST")
	if host == "" {
		host = DefaultServerHost
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultServerPort
	}

	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = DefaultStorageType
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		databaseUrl = DefaultDatabaseUrl
	}

	return Env{
		ServerHost:  host,
		ServerPort:  port,
		StorageType: storageType,
		DatabaseUrl: databaseUrl,
	}
}
