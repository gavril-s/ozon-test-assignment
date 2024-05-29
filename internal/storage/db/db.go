package db

import (
	dbModel "ozon-test-assignment/internal/storage/db/model"
	"ozon-test-assignment/internal/storage/errors"
	storageModel "ozon-test-assignment/internal/storage/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	con *gorm.DB
}

func NewDB(databaseUrl string) (*DB, error) {
	con, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
	if err != nil {
		return nil, errors.DatabaseConnectionFailure{}
	}
	con.AutoMigrate(&dbModel.Post{}, &dbModel.Comment{}, &dbModel.Reply{})
	return &DB{con}, nil
}

func (db *DB) Close() {
	sqlDB, _ := db.con.DB()
	sqlDB.Close()
}

func (db *DB) dbCommentToStorageComment(dbComment *dbModel.Comment) *storageModel.Comment {
	comment := storageModel.Comment{
		ID:          dbComment.ID,
		PostID:      dbComment.PostID,
		ParentID:    dbComment.ParentID,
		Author:      dbComment.Author,
		PublishedAt: dbComment.PublishedAt,
		Content:     dbComment.Content,
		Replies:     make([]*storageModel.Comment, 0),
	}
	return &comment
}

func (db *DB) AddPost(post storageModel.Post) (*storageModel.Post, error) {
	db.con.Create(&post)
	return &post, nil
}

func (db *DB) AddComment(comment storageModel.Comment) (*storageModel.Comment, error) {
	post := struct {
		CommentsEnabled bool
	}{}
	tx := db.con.Table("posts").Where("id = ?", comment.PostID).First(&post)
	if tx.Error != nil {
		return nil, errors.DatabaseQueryExecutionFailure{}
	}
	if !post.CommentsEnabled {
		return nil, errors.CommentsDisabled{}
	}

	dbComment := dbModel.Comment{
		PostID:   comment.PostID,
		ParentID: comment.ParentID,
		Author:   comment.Author,
		Content:  comment.Content,
	}

	tx = db.con.Create(&dbComment)
	if tx.Error != nil {
		return nil, errors.DatabaseQueryExecutionFailure{}
	}
	if comment.ParentID != nil {
		reply := dbModel.Reply{
			PostID:    dbComment.PostID,
			CommentID: *dbComment.ParentID,
			ReplyID:   dbComment.ID,
		}
		tx = db.con.Create(&reply)
		if tx.Error != nil {
			return nil, errors.DatabaseQueryExecutionFailure{}
		}
	}

	return db.dbCommentToStorageComment(&dbComment), nil
}

func (db *DB) GetPost(postId int) (*storageModel.Post, error) {
	var dbPost dbModel.Post
	tx := db.con.First(&dbPost, postId)
	if tx.Error != nil {
		return nil, errors.DatabaseQueryExecutionFailure{}
	}

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
	comments := make(map[int]*storageModel.Comment)
	topLevelComments := make([]*storageModel.Comment, 0, threadLimit)

	dbTopLevelComments := make([]dbModel.Comment, 0, threadLimit)
	tx := db.con.Limit(threadLimit+1).Where("post_id = ?", postId)
	if parentId != nil {
		tx = tx.Where("parent_id = ?", *parentId)
	} else {
		tx = tx.Where("parent_id IS NULL")
	}
	tx = tx.Find(&dbTopLevelComments)
	if tx.Error != nil {
		return nil, errors.DatabaseQueryExecutionFailure{}
	}

	parentIds := make([]int, 0, len(dbTopLevelComments))
	for _, dbComment := range dbTopLevelComments {
		parentIds = append(parentIds, dbComment.ID)
		comment := db.dbCommentToStorageComment(&dbComment)
		topLevelComments = append(topLevelComments, comment)
		comments[comment.ID] = comment
	}

	for depth := 0; depth <= depthLimit; depth++ {
		dbReplies := make([]dbModel.Comment, 0, threadLimit)
		tx = db.con.Raw(
			"SELECT c.id, c.post_id, c.parent_id, c.author, c.published_at, c.content "+
				"FROM replies AS r LEFT JOIN comments AS c ON r.reply_id = c.id "+
				"WHERE r.comment_id IN ?", parentIds,
		).Find(&dbReplies)
		if tx.Error != nil {
			return nil, errors.DatabaseQueryExecutionFailure{}
		}

		parentIds = make([]int, 0, len(dbReplies))
		for _, dbReply := range dbReplies {
			reply := db.dbCommentToStorageComment(&dbReply)
			parentIds = append(parentIds, reply.ID)
			comments[reply.ID] = reply
			if reply.ParentID != nil {
				if parent, exists := comments[*reply.ParentID]; exists {
					parent.Replies = append(parent.Replies, reply)
				}
			}
		}
	}

	return topLevelComments, nil
}

func (db *DB) GetPostsSnippets(snippetLength int, limit int, after *int) ([]*storageModel.PostSnippet, error) {
	snippets := make([]*storageModel.PostSnippet, 0, limit)
	tx := db.con.Table("posts").Select(
		"id AS post_id, title, author, published_at, LEFT(content, ?) AS content_snippet, comments_enabled",
		snippetLength,
	)
	if after != nil {
		tx = tx.Where("id > ?", *after)
	}
	tx = tx.Limit(limit + 1).Find(&snippets)
	if tx.Error != nil {
		return nil, errors.DatabaseQueryExecutionFailure{}
	}
	return snippets, nil
}
