package model

import "time"

type Post struct {
	ID              int
	Title           string
	Author          string
	PublishedAt     time.Time
	Content         string
	CommentsEnabled bool
}

type Comment struct {
	ID          int
	PostID      int
	ParentID    *int
	Author      string
	PublishedAt time.Time
	Content     string
	Replies     []*Comment
}

type PostSnippet struct {
	PostID         int
	Title          string
	Author         string
	PublishedAt    time.Time
	ContentSnippet string
	CommentEnabled bool
}
