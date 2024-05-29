package comment

import (
	graphModel "ozon-test-assignment/graph/model"
	"testing"
	"time"
)

func TestSubscriptions(t *testing.T) {
	comment := &graphModel.Comment{
		ID:     1,
		PostID: 1,
		ContentMeta: &graphModel.ContentMeta{
			Author:      "author",
			PublishedAt: time.Now().String(),
		},
		Content: "Comment sample text",
		Replies: &graphModel.CommentConnection{
			Edges: make([]*graphModel.CommentEdge, 0),
			PageInfo: &graphModel.PageInfo{
				EndCursor:   0,
				HasNextPage: false,
			},
		},
	}

	manager := NewSubscriptionManager()
	id, ch := manager.AddSubscriber(1)
	manager.BroadcastComment(1, comment)

	_, ok := <-ch
	if !ok {
		t.Errorf("Expected non-empty channel")
	}

	manager.RemoveSubscriber(1, id)
	comment.ID = 2
	manager.BroadcastComment(1, comment)

	_, ok = <-ch
	if ok {
		t.Errorf("Channel excpected to be closed!")
	}
}
