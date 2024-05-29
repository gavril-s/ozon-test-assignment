package graph

import (
	"ozon-test-assignment/internal/comment"
	"ozon-test-assignment/internal/storage"
)

type Resolver struct {
	storage                    storage.Storage
	commentSubscriptionManager *comment.SubscriptionManager
}

func NewResolver(storage storage.Storage, commentSubscriptionManager *comment.SubscriptionManager) *Resolver {
	return &Resolver{storage: storage, commentSubscriptionManager: commentSubscriptionManager}
}
