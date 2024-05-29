package db

import (
	"fmt"
	dbModel "ozon-test-assignment/internal/storage/db/model"
	storageModel "ozon-test-assignment/internal/storage/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	con *gorm.DB
}

func NewDB(user, password, host, port, dbname string) (*DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=Europe/Moscow",
		user, password, host, port, dbname,
	)
	var con *gorm.DB
	var err error
	for con == nil || err != nil {
		con, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("datatbase connection failure")
		}
	}
	con.AutoMigrate(&dbModel.Post{}, &dbModel.Comment{}, &dbModel.Reply{})
	return &DB{con}, nil
}

func (db *DB) Close() {
	sqlDB, err := db.con.DB()
	if err != nil {

	}
	sqlDB.Close()
}

func (db *DB) AddPost(post storageModel.Post) (*storageModel.Post, error) {
	db.con.Create(&post)
	return &post, nil
}

func (db *DB) AddComment(comment storageModel.Comment) (*storageModel.Comment, error) {
	post := struct {
		CommentsEnabled bool
	}{}
	db.con.Table("posts").Where("id = ?", comment.PostID).First(&post)
	if !post.CommentsEnabled {
		return nil, fmt.Errorf("comments not allowed")
	}

	dbComment := dbModel.Comment{
		PostID:   comment.PostID,
		ParentID: comment.ParentID,
		Author:   comment.Author,
		Content:  comment.Content,
	}

	db.con.Create(&dbComment)
	if comment.ParentID != nil {
		reply := dbModel.Reply{
			PostID:    dbComment.PostID,
			CommentID: *dbComment.ParentID,
			ReplyID:   dbComment.ID,
		}
		db.con.Create(&reply)
	}

	comment = storageModel.Comment{
		ID:          dbComment.ID,
		PostID:      dbComment.PostID,
		ParentID:    dbComment.ParentID,
		Author:      dbComment.Author,
		PublishedAt: dbComment.PublishedAt,
		Content:     dbComment.Content,
		Replies:     make([]*storageModel.Comment, 0),
	}
	return &comment, nil
}

func (db *DB) GetPost(postId int) (*storageModel.Post, error) {
	var dbPost dbModel.Post
	db.con.First(&dbPost, postId)

	post := storageModel.Post{
		ID:              dbPost.ID,
		Title:           dbPost.Title,
		Author:          dbPost.Author,
		PublishedAt:     dbPost.PublishedAt,
		Content:         dbPost.Content,
		CommentsEnabled: dbPost.CommentsEnabled,
	}
	return &post, nil
}

func (db *DB) GetComments(postId int, parentId *int, depthLimit int, threadLimit int, after *int) ([]*storageModel.Comment, error) {
	topLevelComments := make([]*storageModel.Comment, 0, threadLimit)
	comments := make(map[int]*storageModel.Comment)

	dbTopLevelComments := make([]dbModel.Comment, 0, threadLimit)
	tx := db.con.Order("id").Limit(threadLimit)
	if parentId != nil {
		tx = tx.Where("comment_id = ?", *parentId)
	}
	tx = tx.Find(&dbTopLevelComments)
	if tx.Error != nil {
		return nil, tx.Error
	}
	parentIds := make([]int, 0, len(dbTopLevelComments))
	for _, dbComment := range dbTopLevelComments {
		parentIds = append(parentIds, dbComment.ID)
		comment := storageModel.Comment{
			ID:          dbComment.ID,
			PostID:      dbComment.PostID,
			ParentID:    dbComment.ParentID,
			Author:      dbComment.Author,
			PublishedAt: dbComment.PublishedAt,
			Content:     dbComment.Content,
			Replies:     make([]*storageModel.Comment, 0),
		}
		topLevelComments = append(topLevelComments, &comment)
		comments[comment.ID] = &comment
	}
	for depth := 1; depth <= depthLimit; depth++ {
		dbReplies := make([]dbModel.Comment, 0, threadLimit)
		db.con.Raw(
			"SELECT c.id, c.post_id, c.parent_id, c.author, c.published_at, c.content "+
				"FROM replies AS r LEFT JOIN comments AS c ON r.reply_id = c.id "+
				"WHERE r.comment_id IN ?", parentIds,
		).Find(&dbReplies)
		if tx.Error != nil {
			return nil, tx.Error
		}
		parentIds = make([]int, 0, len(dbReplies))
		for _, dbReply := range dbReplies {
			reply := storageModel.Comment{
				ID:          dbReply.ID,
				PostID:      dbReply.PostID,
				ParentID:    dbReply.ParentID,
				Author:      dbReply.Author,
				PublishedAt: dbReply.PublishedAt,
				Content:     dbReply.Content,
				Replies:     make([]*storageModel.Comment, 0),
			}
			parentIds = append(parentIds, reply.ID)
			comments[reply.ID] = &reply
			if reply.ParentID != nil {
				if parent, exists := comments[*reply.ParentID]; exists {
					parent.Replies = append(parent.Replies, &reply)
				}
			}
		}
	}
	return topLevelComments, nil
}

func (db *DB) GetPostsSnippets(snippetLength int, limit int, after *int) ([]*storageModel.PostSnippet, error) {
	snippets := make([]*storageModel.PostSnippet, 0, limit)
	tx := db.con.Table("posts").Select(
		"id, title, author, published_at, LEFT(content, ?), comments_enabled",
		snippetLength,
	)
	if after != nil {
		tx = tx.Where("id > ?", *after)
	}
	tx = tx.Limit(limit).Find(&snippets)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return snippets, nil
}
