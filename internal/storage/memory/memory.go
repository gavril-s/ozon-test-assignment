package memory

import (
	"ozon-test-assignment/internal/storage/errors"
	storageModel "ozon-test-assignment/internal/storage/model"
	"sync"
)

type Memory struct {
	mu               sync.Mutex
	posts            []*storageModel.Post
	comments         []*storageModel.Comment
	topLevelComments map[int][]*storageModel.Comment
}

func NewMemory() *Memory {
	posts := make([]*storageModel.Post, 0)
	comments := make([]*storageModel.Comment, 0)
	topLevelComments := make(map[int][]*storageModel.Comment)
	return &Memory{
		posts:            posts,
		comments:         comments,
		topLevelComments: topLevelComments,
	}
}

func (m *Memory) containsPost(postId int) bool {
	return postId >= 0 && postId <= len(m.posts)
}

func (m *Memory) containsComment(commentId int) bool {
	return commentId >= 0 && commentId <= len(m.comments)
}

func (m *Memory) AddPost(post storageModel.Post) (*storageModel.Post, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	post.ID = len(m.posts) + 1
	m.posts = append(m.posts, &post)
	m.topLevelComments[post.ID] = make([]*storageModel.Comment, 0)
	return &post, nil
}

func (m *Memory) AddComment(comment storageModel.Comment) (*storageModel.Comment, error) {
	if !m.containsPost(comment.PostID) {
		return nil, errors.PostNotFoundError{}
	}
	post := m.posts[comment.PostID-1]
	if !post.CommentsEnabled {
		return nil, errors.CommentsDisabled{}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	comment.ID = len(m.comments) + 1
	comment.Replies = make([]*storageModel.Comment, 0)
	if comment.ParentID != nil {
		parentID := *comment.ParentID
		comment.ParentID = &parentID
	}
	m.comments = append(m.comments, &comment)

	if comment.ParentID == nil {
		m.topLevelComments[comment.PostID] = append(m.topLevelComments[comment.PostID], &comment)
	} else {
		if !m.containsComment(*comment.ParentID) {
			return nil, errors.ParentCommentNotFound{}
		}
		parentComment := m.comments[*comment.ParentID-1]
		parentComment.Replies = append(parentComment.Replies, &comment)
		m.comments[*comment.ParentID-1] = parentComment
	}

	return &comment, nil
}

func (m *Memory) GetPost(postId int) (*storageModel.Post, error) {
	if !m.containsPost(postId) {
		return nil, errors.PostNotFoundError{}
	}
	post := m.posts[postId-1]
	return post, nil
}

func (m *Memory) GetComments(postId int, parentId *int, depthLimit int, threadLimit int, after *int) ([]*storageModel.Comment, error) {
	if !m.containsPost(postId) {
		return nil, errors.PostNotFoundError{}
	}

	var parentComments []*storageModel.Comment
	if parentId == nil {
		parentComments = m.topLevelComments[postId]
	} else {
		if !m.containsComment(*parentId) {
			return nil, errors.ParentCommentNotFound{}
		}
		parent := m.comments[*parentId-1]
		parentComments = parent.Replies
	}

	comments := []*storageModel.Comment{}
	for _, comment := range parentComments {
		if after != nil && comment.ID <= *after {
			continue
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (m *Memory) GetPostsSnippets(snippetLength int, limit int, after *int) ([]*storageModel.PostSnippet, error) {
	start := 1
	if after != nil {
		start = *after + 1
	}
	snippets := make([]*storageModel.PostSnippet, 0, limit)
	for id := start; id <= len(m.posts); id++ {
		if len(snippets) > limit {
			break
		}
		post := m.posts[id-1]
		contentSnippet := post.Content
		if len(contentSnippet) >= snippetLength {
			contentSnippet = contentSnippet[:snippetLength]
		}
		snippet := storageModel.PostSnippet{
			PostID:         post.ID,
			Title:          post.Title,
			Author:         post.Author,
			PublishedAt:    post.PublishedAt,
			ContentSnippet: contentSnippet,
		}
		snippets = append(snippets, &snippet)
	}
	return snippets, nil
}
