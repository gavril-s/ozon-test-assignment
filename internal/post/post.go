package post

import (
	"ozon-test-assignment/configs"
	graphModel "ozon-test-assignment/graph/model"
	"ozon-test-assignment/internal/storage"
	storageModel "ozon-test-assignment/internal/storage/model"
)

func createPostSnippetConnection(
	snippets []*storageModel.PostSnippet,
	limit int,
) *graphModel.PostSnippetConnection {
	edges := make([]*graphModel.PostSnippetEdge, 0, limit)
	hasNextPage := false
	endCursor := 0

	for _, snippet := range snippets {
		if len(edges) >= limit {
			hasNextPage = true
			break
		}
		edge := graphModel.PostSnippetEdge{
			Cursor: snippet.PostID,
			Node: &graphModel.PostSnippet{
				PostID: snippet.PostID,
				Title:  snippet.Title,
				ContentMeta: &graphModel.ContentMeta{
					Author:      snippet.Author,
					PublishedAt: snippet.PublishedAt.String(),
				},
				ContentSnippet: snippet.ContentSnippet,
			},
		}
		edges = append(edges, &edge)
		endCursor = edge.Cursor
	}

	return &graphModel.PostSnippetConnection{
		Edges: edges,
		PageInfo: &graphModel.PageInfo{
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		},
	}
}

func Post(storage storage.Storage, id int) (*graphModel.Post, error) {
	dbPost, err := storage.GetPost(id)
	if err != nil {
		return nil, err
	}

	post := graphModel.Post{
		ID:    dbPost.ID,
		Title: dbPost.Title,
		ContentMeta: &graphModel.ContentMeta{
			Author:      dbPost.Author,
			PublishedAt: dbPost.PublishedAt.String(),
		},
		Content:         dbPost.Content,
		CommentsEnabled: dbPost.CommentsEnabled,
	}
	return &post, nil
}

func PostSnippets(storage storage.Storage, snippetLength *int, first *int, after *int) (*graphModel.PostSnippetConnection, error) {
	limit := configs.PostSnippetsLimit
	if first != nil {
		limit = *first
	}
	length := configs.PostSnippetDefaultLength
	if snippetLength != nil {
		length = *snippetLength
	}
	snippets, err := storage.GetPostsSnippets(length, limit, after)
	if err != nil {
		return nil, err
	}
	connection := createPostSnippetConnection(snippets, limit)
	return connection, nil
}

func CreatePost(storage storage.Storage, title string, author string, content string, commentsEnabled bool) (*graphModel.Post, error) {
	dbPost := &storageModel.Post{
		Title:           title,
		Author:          author,
		Content:         content,
		CommentsEnabled: commentsEnabled,
	}

	dbPost, err := storage.AddPost(*dbPost)
	if err != nil {
		return nil, err
	}

	post := graphModel.Post{
		ID:    dbPost.ID,
		Title: dbPost.Title,
		ContentMeta: &graphModel.ContentMeta{
			Author:      dbPost.Author,
			PublishedAt: dbPost.PublishedAt.String(),
		},
		Content:         dbPost.Content,
		CommentsEnabled: dbPost.CommentsEnabled,
	}
	return &post, nil
}
