package expandenv

import (
	"bytes"
	"os"
)

type Expr interface {
	Eval(Writer) (int, error)
}

type Writer interface {
	WriteString(string) (int, error)
	WriteRune(rune) (int, error)
}

type List []Expr

func (list List) Eval(w Writer) (n int, err error) {
	for _, expr := range list {
		var m int
		m, err = expr.Eval(w)
		n += m
		if err != nil {
			break
		}
	}
	return
}

type Ident string

func (ident Ident) Eval(w Writer) (n int, err error) {
	return w.WriteString(os.Getenv(string(ident)))
}

type Rune rune

func (r Rune) Eval(w Writer) (n int, err error) {
	return w.WriteRune(rune(r))
}

type Elvis struct {
	ident    Ident
	op       rune
	fallback Expr
}

func (elvis *Elvis) Eval(w Writer) (n int, err error) {
	key := string(elvis.ident)
	if v, ok := os.LookupEnv(key); ok {
		return w.WriteString(v)
	}
	if elvis.op == '=' {
		b := new(bytes.Buffer)
		if _, err := elvis.fallback.Eval(b); err != nil {
			return 0, err
		}
		s := b.String()
		if err := os.Setenv(key, s); err != nil {
			return 0, err
		}
		return w.WriteString(s)
	}
	return elvis.fallback.Eval(w)
}
