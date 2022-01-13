package db

//DuplicateKeyError used when occurs a constraint violation of unique key
type DuplicateKeyError struct{}

func (e *DuplicateKeyError) Error() string {
	return "duplicated beer id"
}

//NotFoundError used when no rows were found in a select query
type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "beer not found"
}
