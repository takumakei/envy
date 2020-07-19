package expandenv

import (
	"bytes"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test1(t *testing.T) {
	k, v := "ABCDEFG", "abcdefg"
	if err := os.Setenv(k, v); err != nil {
		t.Fatal(err)
	}

	s := "$" + k

	var p Parser
	expr, _ := p.Init(s).Parse()
	b := new(bytes.Buffer)
	t.Log(spew.Sdump(expr))
	expr.Eval(b)
	if r := b.String(); r != v {
		t.Fatalf("%q != %q", v, r)
	}
}
