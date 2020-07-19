package expandenv

import (
	"bytes"
)

// ExpandEnv replaces ${var} or $var in the string according to the values
// of the current environment variables. References to undefined
// variables are replaced by the empty string.
//
// ExpandEnv replaces these expressions.
//
//   ${var:-expr}
//   ${var:=expr}
//
// In case of the variables var is undefined, ExpandEnv replaces it by expr.
// In addition to that, ':=' set the environment variable var to the expr.
func ExpandEnv(s string) string {
	var p Parser
	expr, _ := p.Init(s).Parse()
	b := new(bytes.Buffer)
	expr.Eval(b)
	return b.String()
}
