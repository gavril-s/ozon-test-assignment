package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"
	"ozon-test-assignment/graph/model"
	"ozon-test-assignment/internal/comment"
	"ozon-test-assignment/internal/post"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, title string, author string, content string, commentsEnabled bool) (*model.Post, error) {
	return post.CreatePost(r.storage, title, author, content, commentsEnabled)
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, postID int, parentID *int, author string, content string) (*model.Comment, error) {
	comment, err := comment.CreateComment(r.storage, postID, parentID, author, content)
	r.commentSubscriptionManager.BroadcastComment(postID, comment)
	return comment, err
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id int) (*model.Post, error) {
	return post.Post(r.storage, id)
}

// PostSnippets is the resolver for the postSnippets field.
func (r *queryResolver) PostSnippets(ctx context.Context, snippetLength *int, first *int, after *int) (*model.PostSnippetConnection, error) {
	return post.PostSnippets(r.storage, snippetLength, first, after)
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, postID int, parentID *int, first *int, depth *int, after *int) (*model.CommentConnection, error) {
	comments, err := comment.Comments(r.storage, postID, parentID, first, depth, after)
	return comments, err
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID int) (<-chan *model.Comment, error) {
	id, commentChan := r.commentSubscriptionManager.AddSubscriber(postID)
	go func() {
		<-ctx.Done()
		r.commentSubscriptionManager.RemoveSubscriber(postID, id)
	}()
	return commentChan, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
