package msgrelay

import "errors"

var (
	ErrDuplicateSession = errors.New(
		"user already has this session")
)
