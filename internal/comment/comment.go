package comment

import (
	"fmt"
	"ozon-test-assignment/graph/model"
	graphModel "ozon-test-assignment/graph/model"
	"ozon-test-assignment/internal/storage"
	storageModel "ozon-test-assignment/internal/storage/model"
)

const CommentsDepthLimit = 5
const CommentsThreadLimit = 100
const ContentLengthLimit = 2000

func createCommentConnection(comments []*storageModel.Comment, depthLimit int, threadLimit int) *graphModel.CommentConnection {
	edges := make([]*graphModel.CommentEdge, 0, threadLimit)
	hasNextPage := false
	endCursor := 0
	for _, comment := range comments {
		if len(edges) >= threadLimit {
			hasNextPage = true
			break
		}
		replies := &graphModel.CommentConnection{}
		if depthLimit > 0 {
			replies = createCommentConnection(comment.Replies, depthLimit-1, threadLimit)
		}
		if depthLimit == 1 {
			replies = &graphModel.CommentConnection{
				Edges: []*model.CommentEdge{},
				PageInfo: &graphModel.PageInfo{
					EndCursor:   0,
					HasNextPage: replies.PageInfo.HasNextPage || len(replies.Edges) > 0,
				},
			}
		}
		edge := model.CommentEdge{
			Cursor: comment.ID,
			Node: &graphModel.Comment{
				ID:     comment.ID,
				PostID: comment.PostID,
				ContentMeta: &graphModel.ContentMeta{
					Author:      comment.Author,
					PublishedAt: comment.PublishedAt.String(),
				},
				Content: comment.Content,
				Replies: replies,
			},
		}
		edges = append(edges, &edge)
		endCursor = edge.Cursor
	}
	return &graphModel.CommentConnection{
		Edges: edges,
		PageInfo: &graphModel.PageInfo{
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		},
	}
}

func Comments(
	storage storage.Storage,
	postID int,
	parentID *int,
	first *int,
	depth *int,
	after *int,
) (*graphModel.CommentConnection, error) {
	depthLimit := CommentsDepthLimit
	if depth != nil {
		depthLimit = *depth
	}
	threadLimit := CommentsThreadLimit
	if first != nil {
		threadLimit = *first
	}
	comments, err := storage.GetComments(postID, parentID, depthLimit, threadLimit, after)
	if err != nil {
		return nil, fmt.Errorf("get comment error")
	}
	connection := createCommentConnection(comments, depthLimit, threadLimit)
	return connection, nil
}

func CreateComment(
	storage storage.Storage,
	postID int,
	parentID *int,
	author string,
	content string,
) (*graphModel.Comment, error) {
	if len(content) > ContentLengthLimit {
		return nil, fmt.Errorf("limit comment")
	}

	commentData := storageModel.Comment{
		PostID:   postID,
		ParentID: parentID,
		Author:   author,
		Content:  content,
	}
	comment, err := storage.AddComment(commentData)
	if err != nil {
		return nil, err
	}

	graphComment := graphModel.Comment{
		ID:     comment.ID,
		PostID: comment.PostID,
		ContentMeta: &graphModel.ContentMeta{
			Author:      comment.Author,
			PublishedAt: comment.PublishedAt.String(),
		},
		Content: comment.Content,
		Replies: &graphModel.CommentConnection{
			Edges: make([]*graphModel.CommentEdge, 0),
			PageInfo: &graphModel.PageInfo{
				EndCursor:   0,
				HasNextPage: false,
			},
		},
	}
	return &graphComment, nil
}
