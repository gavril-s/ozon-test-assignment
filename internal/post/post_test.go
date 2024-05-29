package post

import (
	"ozon-test-assignment/internal/storage"
	"ozon-test-assignment/internal/storage/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestStorage() storage.Storage {
	testStorage := memory.NewMemory()
	return testStorage
}

func TestPosts(t *testing.T) {
	testStorage := getTestStorage()
	CreatePost(testStorage, "post 1", "author", "post 1 other text", false)
	CreatePost(testStorage, "post 2", "author", "post 2 other text", false)
	CreatePost(testStorage, "post 3", "author", "post 3 other text", false)

	post, err := Post(testStorage, 1)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "post 1", post.Title)

	length := 4
	snippets, err := PostSnippets(testStorage, &length, nil, nil)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 3, len(snippets.Edges))
	for _, edge := range snippets.Edges {
		assert.Equal(t, "post", edge.Node.ContentSnippet)
	}
}
