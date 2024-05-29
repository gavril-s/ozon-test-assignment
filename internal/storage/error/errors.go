package error

type PostNotFoundError struct{}

func (e *PostNotFoundError) Error() string {
	return "Post Not Found"
}
