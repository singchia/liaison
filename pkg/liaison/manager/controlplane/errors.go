package controlplane

import "errors"

// errForbidden is returned when a request targets a resource owned by
// another user. HTTP handlers map this to 403.
var errForbidden = errors.New("forbidden")

// ErrForbidden returns the sentinel used for cross-user access denials, for
// callers outside this package that need errors.Is matching.
func ErrForbidden() error { return errForbidden }
