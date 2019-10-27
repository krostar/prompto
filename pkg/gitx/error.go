package gitx

// ErrRepositoryDoesNotExists means repository has not been found.
const ErrRepositoryDoesNotExists gitError = "repository does not exists"

type gitError string

func (e gitError) Error() string { return string(e) }
