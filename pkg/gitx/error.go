package gitx

const ErrRepositoryDoesNotExists gitError = "repository does not exists"

type gitError string

func (e gitError) Error() string {
	return string(e)
}
