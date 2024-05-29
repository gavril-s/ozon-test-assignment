package post

import (
	"fmt"
	graphModel "ozon-test-assignment/graph/model"
	"ozon-test-assignment/internal/storage"
	storageModel "ozon-test-assignment/internal/storage/model"
)

const SnippetsLimit = 100
const SnippetDefaultLength = 300

func Post(storage storage.Storage, id int) (*graphModel.Post, error) {
	post, err := storage.GetPost(id)
	if err != nil {
		return nil, fmt.Errorf("get post error")
	}

	graphPost := graphModel.Post{
		ID:    post.ID,
		Title: post.Title,
		ContentMeta: &graphModel.ContentMeta{
			Author:      post.Author,
			PublishedAt: post.PublishedAt.String(),
		},
		Content:         post.Content,
		CommentsEnabled: post.CommentsEnabled,
	}
	return &graphPost, nil
}

func PostSnippets(storage storage.Storage, snippetLength *int, first *int, after *int) (*graphModel.PostSnippetConnection, error) {
	limit := SnippetsLimit
	if first != nil {
		limit = *first
	}

	length := SnippetDefaultLength
	if snippetLength != nil {
		length = *snippetLength
	}

	snippets, err := storage.GetPostsSnippets(length, limit, after)
	if err != nil {
		return nil, fmt.Errorf("get snippets error")
	}

	var snippetEdges []*graphModel.PostSnippetEdge
	for _, snippet := range snippets {
		graphSnippet := graphModel.PostSnippet{
			PostID: snippet.PostID,
			Title:  snippet.Title,
			ContentMeta: &graphModel.ContentMeta{
				Author:      snippet.Author,
				PublishedAt: snippet.PublishedAt.String(),
			},
			ContentSnippet: snippet.ContentSnippet,
		}
		snippetEdges = append(snippetEdges, &graphModel.PostSnippetEdge{
			Cursor: graphSnippet.PostID,
			Node:   &graphSnippet,
		})
	}

	endCursor := 0
	hasNextPage := false
	if len(snippets) > 0 {
		endCursor = snippets[len(snippets)-1].PostID
		hasNextPage = len(snippets) == limit
	}

	connection := &graphModel.PostSnippetConnection{
		Edges: snippetEdges,
		PageInfo: &graphModel.PageInfo{
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		},
	}
	return connection, nil
}

func CreatePost(storage storage.Storage, title string, author string, content string, commentsEnabled bool) (*graphModel.Post, error) {
	postData := storageModel.Post{
		Title:           title,
		Author:          author,
		Content:         content,
		CommentsEnabled: commentsEnabled,
	}
	post, err := storage.AddPost(postData)
	if err != nil {
		return nil, fmt.Errorf("create post error")
	}

	graphPost := graphModel.Post{
		ID:    post.ID,
		Title: post.Title,
		ContentMeta: &graphModel.ContentMeta{
			Author:      post.Author,
			PublishedAt: post.PublishedAt.String(),
		},
		Content:         post.Content,
		CommentsEnabled: post.CommentsEnabled,
	}
	return &graphPost, nil
}
