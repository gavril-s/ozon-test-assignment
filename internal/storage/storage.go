package storage

import "ozon-test-assignment/internal/storage/model"

const (
	StorageTypeDB     = "DB"
	StorageTypeMemory = "MEMORY"
)

type Storage interface {
	AddPost(post model.Post) (*model.Post, error)
	AddComment(comment model.Comment) (*model.Comment, error)
	GetPost(postId int) (*model.Post, error)
	GetComments(postId int, parentId *int, depthLimit int, threadLimit int, after *int) ([]*model.Comment, error)
	GetPostsSnippets(snippetLength int, limit int, after *int) ([]*model.PostSnippet, error)
}
