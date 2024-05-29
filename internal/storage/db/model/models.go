package model

import "time"

type Post struct {
	ID              int       `gorm:"type:uint;primaryKey;autoIncrement;not null"`
	Title           string    `gorm:"not null"`
	Author          string    `gorm:"not null"`
	PublishedAt     time.Time `gorm:"autoCreateTime;not null"`
	Content         string    `gorm:"not null"`
	CommentsEnabled bool      `gorm:"not null"`
}

type Comment struct {
	ID          int       `gorm:"type:uint;primaryKey;autoIncrement;not null"`
	PostID      int       `gorm:"type:uint;not null"`
	ParentID    *int      `gorm:"type:uint"`
	Author      string    `gorm:"not null"`
	PublishedAt time.Time `gorm:"not null"`
	Content     string    `gorm:"size:2000;not null"`
}

type Reply struct {
	ReplyID   int `gorm:"type:uint;primaryKey;not null"`
	CommentID int `gorm:"type:uint;not null"`
	PostID    int `gorm:"type:uint;not null"`
}
