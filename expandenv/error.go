package expandenv

import "errors"

var (
	UnexpectedEOF  = errors.New("missing '}', unexpected EOF")
	UnexpectedRune = errors.New("must be '-' or '='")
)
