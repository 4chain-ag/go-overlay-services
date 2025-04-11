package core

import (
	"errors"
)

var ErrGASPVersionMismatch = errors.New("gasp version mismatch")

func (e *GASPVersionMismatchError) Is(target error) bool {
	_, ok := target.(*GASPVersionMismatchError)
	return ok
}
