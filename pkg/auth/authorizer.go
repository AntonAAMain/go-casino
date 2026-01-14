package auth

type Authorizer interface {
	Authorize(token string) (userID uint, role string, err error)
}
