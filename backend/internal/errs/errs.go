package errs

import "errors"

var (
	NotFound = errors.New("not found")
	Empty = errors.New("the repo is empty")
	Internal = errors.New("internal error")
	InvalidPassword = errors.New("invalid password")
	InvalidLogin = errors.New("invalid login")
	RefreshExpired  = errors.New("refresh token expired")
)
