package comment

import (
	"fmt"
	"ozon-test-assignment/configs"
	graphModel "ozon-test-assignment/graph/model"
	"ozon-test-assignment/internal/storage"
	storageModel "ozon-test-assignment/internal/storage/model"
)

func createCommentConnection(
	comments []*storageModel.Comment,
	depthLimit int,
	threadLimit int,
) *graphModel.CommentConnection {
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
				Edges: []*graphModel.CommentEdge{},
				PageInfo: &graphModel.PageInfo{
					EndCursor:   0,
					HasNextPage: replies.PageInfo.HasNextPage || len(replies.Edges) > 0,
				},
			}
		}

		edge := graphModel.CommentEdge{
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
	depthLimit := configs.CommentsDepthLimit
	if depth != nil {
		depthLimit = *depth
	}
	threadLimit := configs.CommentsThreadLimit
	if first != nil {
		threadLimit = *first
	}
	comments, err := storage.GetComments(postID, parentID, depthLimit, threadLimit, after)
	if err != nil {
		return nil, err
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
	if len(content) > configs.CommentsContentLengthLimit {
		return nil, fmt.Errorf("Error")
	}

	dbComment := &storageModel.Comment{
		PostID:   postID,
		ParentID: parentID,
		Author:   author,
		Content:  content,
	}

	dbComment, err := storage.AddComment(*dbComment)
	if err != nil {
		return nil, err
	}

	comment := graphModel.Comment{
		ID:     dbComment.ID,
		PostID: dbComment.PostID,
		ContentMeta: &graphModel.ContentMeta{
			Author:      dbComment.Author,
			PublishedAt: dbComment.PublishedAt.String(),
		},
		Content: dbComment.Content,
		Replies: &graphModel.CommentConnection{
			Edges: make([]*graphModel.CommentEdge, 0),
			PageInfo: &graphModel.PageInfo{
				EndCursor:   0,
				HasNextPage: false,
			},
		},
	}
	return &comment, nil
}
