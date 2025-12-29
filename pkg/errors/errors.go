package errors

import (
	"github.com/pkg/errors"
)

var (
	WithStack = errors.WithStack
	Wrap      = errors.Wrap
	Wrapf     = errors.Wrapf
	Is        = errors.Is
	Errorf    = errors.Errorf
)

const (
	DefaultBadRequestID   = "bad_request"
	DefaultUnauthorizedID = "unauthorized"
	DefaultForbiddenID    = "forbidden"
	DefaultNotFoundID     = "not_found"
)

type Error struct {
	ID     string
	Code   int32
	Detail string
	Status string
}
