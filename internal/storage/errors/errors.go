package errors

import "fmt"

type PostNotFoundError struct{}

func (e PostNotFoundError) Error() string {
	return "Post Not Found"
}

type CommentsDisabled struct{}

func (e CommentsDisabled) Error() string {
	return "Comments Disabled"
}

type ParentCommentNotFound struct{}

func (e ParentCommentNotFound) Error() string {
	return "Parent Comment Not Found"
}

type DatabaseConnectionFailure struct{}

func (e DatabaseConnectionFailure) Error() string {
	return "Database Connection Failure"
}

type DatabaseQueryExecutionFailure struct {
	Query *string
}

func (e DatabaseQueryExecutionFailure) Error() string {
	if e.Query != nil {
		return fmt.Sprintf("Database Query Execution Failure (Query: %s)", *e.Query)
	} else {
		return "Database Query Execution Failure"
	}
}
