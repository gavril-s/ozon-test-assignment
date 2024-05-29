package memory

import (
	"fmt"
	"ozon-test-assignment/internal/storage/model"
	"sync"
)

type Memory struct {
	mu               sync.Mutex
	posts            []*model.Post
	comments         map[int]*model.Comment
	topLevelComments map[int][]*model.Comment
	lastCommentId    int
}

func NewMemory() *Memory {
	posts := make([]*model.Post, 0)
	comments := make(map[int]*model.Comment)
	topLevelComments := make(map[int][]*model.Comment)
	return &Memory{
		posts:            posts,
		comments:         comments,
		topLevelComments: topLevelComments,
		lastCommentId:    0,
	}
}

func (m *Memory) containsPost(postId int) bool {
	return postId >= 0 && postId <= len(m.posts)
}

func (m *Memory) AddPost(post model.Post) (*model.Post, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	post.ID = len(m.posts) + 1
	m.posts = append(m.posts, &post)
	m.topLevelComments[post.ID] = make([]*model.Comment, 0)
	return &post, nil
}

func (m *Memory) AddComment(comment model.Comment) (*model.Comment, error) {
	if !m.containsPost(comment.PostID) {
		return nil, fmt.Errorf("post not found")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastCommentId++
	comment.ID = m.lastCommentId
	comment.Replies = make([]*model.Comment, 0)

	m.comments[comment.ID] = &comment
	if comment.ParentID == nil {
		m.topLevelComments[comment.PostID] = append(m.topLevelComments[comment.PostID], &comment)
	} else {
		parentComment, exists := m.comments[*comment.ParentID]
		if !exists {
			return nil, fmt.Errorf("parent id does not exist")
		}
		parentComment.Replies = append(parentComment.Replies, &comment)
		m.comments[*comment.ParentID] = parentComment
	}

	for k, v := range m.comments {
		fmt.Println(k, v)
	}
	return &comment, nil
}

func (m *Memory) GetPost(postId int) (*model.Post, error) {
	if !m.containsPost(postId) {
		return nil, fmt.Errorf("post not found")
	}
	post := m.posts[postId]
	return post, nil
}

func (m *Memory) GetComments(postId int, parentId *int, depthLimit int, threadLimit int, after *int) ([]*model.Comment, error) {
	if !m.containsPost(postId) {
		return nil, fmt.Errorf("post not found")
	}

	var parentComments []*model.Comment
	if parentId == nil {
		parentComments = m.topLevelComments[postId]
	} else {
		parent, exists := m.comments[*parentId]
		if !exists {
			fmt.Errorf("parent doesnt exist")
		}
		parentComments = parent.Replies
	}

	comments := []*model.Comment{}
	for _, comment := range parentComments {
		if after != nil && comment.ID <= *after {
			continue
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (m *Memory) GetPostsSnippets(snippetLength int, limit int, after *int) ([]*model.PostSnippet, error) {
	start := 0
	if after != nil {
		start = *after
	}
	snippets := make([]*model.PostSnippet, 0, limit)
	for id := start; id < len(m.posts); id++ {
		post := m.posts[id]
		contentSnippet := post.Content
		if len(contentSnippet) >= limit {
			contentSnippet = contentSnippet[:limit]
		}
		snippet := model.PostSnippet{
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
