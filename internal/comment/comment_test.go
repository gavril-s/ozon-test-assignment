package comment

import (
	graphModel "ozon-test-assignment/graph/model"
	"ozon-test-assignment/internal/storage"
	"ozon-test-assignment/internal/storage/memory"
	storageModel "ozon-test-assignment/internal/storage/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestStorage() storage.Storage {
	testPost := storageModel.Post{
		Title:           "Title",
		Author:          "Author",
		Content:         "Post content",
		CommentsEnabled: true,
	}
	testComment := storageModel.Comment{
		PostID:  1,
		Author:  "Commenter",
		Content: "Comment content",
	}

	testStorage := memory.NewMemory()
	testStorage.AddPost(testPost)
	testStorage.AddComment(testComment)
	return testStorage
}

func TestCommentsBasic(t *testing.T) {
	testStorage := getTestStorage()
	comments, err := Comments(testStorage, 1, nil, nil, nil, nil)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, len(comments.Edges))
}

func TestCommentsComplex(t *testing.T) {
	testStorage := getTestStorage()
	parentId := 1
	comment := storageModel.Comment{
		PostID:   1,
		ParentID: &parentId,
		Author:   "Commenter",
		Content:  "Reply 1",
	}
	createdFirst, err := testStorage.AddComment(comment)
	if err != nil {
		t.Error(err)
	}

	comment.Content = "Reply 2"
	createdSecond, err := testStorage.AddComment(comment)
	if err != nil {
		t.Error(err)
	}

	parentId = createdFirst.ID
	comment.Content = "Reply 1 to Reply 1"
	_, err = testStorage.AddComment(comment)
	if err != nil {
		t.Error(err)
	}

	comment.Content = "Reply 2 to Reply 1"
	_, err = testStorage.AddComment(comment)
	if err != nil {
		t.Error(err)
	}

	parentId = createdSecond.ID
	comment.Content = "Reply 1 to Reply 2"
	_, err = testStorage.AddComment(comment)
	if err != nil {
		t.Error(err)
	}

	comment.Content = "Reply 2 to Reply 2"
	_, err = testStorage.AddComment(comment)
	if err != nil {
		t.Error(err)
	}

	first := 1
	depth := 2
	parentId = 1
	comments, err := Comments(testStorage, 1, &parentId, &first, &depth, nil)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, true, comments.PageInfo.HasNextPage)
	assert.Equal(t, 1, len(comments.Edges))
	assert.Equal(t, 1, len(comments.Edges[0].Node.Replies.Edges))
	assert.Equal(t, true, comments.Edges[0].Node.Replies.PageInfo.HasNextPage)
}

func TestCreateComment(t *testing.T) {
	testStorage := getTestStorage()
	expected := graphModel.Comment{
		ID:     2,
		PostID: 1,
		ContentMeta: &graphModel.ContentMeta{
			Author:      "Commenter 2",
			PublishedAt: "",
		},
		Content: "Comment content 2",
		Replies: &graphModel.CommentConnection{
			Edges: make([]*graphModel.CommentEdge, 0),
			PageInfo: &graphModel.PageInfo{
				EndCursor:   0,
				HasNextPage: false,
			},
		},
	}

	parentId := 1
	created, err := CreateComment(testStorage, 1, &parentId, "Commenter 2", "Comment content 2")
	if err != nil {
		t.Error(err)
	}
	expected.ContentMeta.PublishedAt = created.ContentMeta.PublishedAt

	assert.Equal(t, *expected.ContentMeta, *created.ContentMeta, "Content meta differ")
	assert.Equal(t, *expected.Replies.PageInfo, *created.Replies.PageInfo, "Replies page info differ")
	assert.Equal(t, expected.Replies.Edges, created.Replies.Edges, "Replise edges differ")
	assert.Equal(t, expected.ID, created.ID)
	assert.Equal(t, expected.PostID, created.PostID)
	assert.Equal(t, expected.Content, created.Content)
}
