package expandenv

import (
	"io"
	"strings"
	"unicode"
)

type Parser struct {
	src  *strings.Reader
	read int
	Pos  int
}

func (p *Parser) Init(s string) *Parser {
	p.src = strings.NewReader(s)
	p.read = 0
	p.Pos = 0
	return p
}

func (p *Parser) readRune() (rune, error) {
	r, n, err := p.src.ReadRune()
	p.read = n
	p.Pos += n
	return r, err
}

func (p *Parser) unreadRune() error {
	if p.read == 0 {
		panic("assert(p.read > 0)")
	}
	p.Pos -= p.read
	p.read = 0
	return p.src.UnreadRune()
}

func (p *Parser) Parse() (Expr, error) {
	var list List
	var err error
	for {
		var r0 rune
		r0, err = p.readRune()
		if err != nil {
			break
		}

		var expr Expr
		expr, err = p.parse(r0)
		if err != nil {
			break
		}

		list = append(list, expr)
	}
	return list, err
}

func (p *Parser) parse(r0 rune) (Expr, error) {
	switch r0 {
	case '$':
		return p.parseExpand(r0)

	case '\\':
		return p.parseEscape(r0)
	}

	return Rune(r0), nil
}

func (p *Parser) parseExpand(r0 rune) (Expr, error) {
	r1, err := p.readRune()
	if err != nil {
		return Rune(r0), err
	}

	switch {
	case r1 == '{':
		return p.parseExpr(r1)

	case IsIdentHead(r1):
		expr, err := p.parseIdent(r1)
		if err == io.EOF {
			err = nil
		}
		return expr, err
	}

	return Rune(r0), p.unreadRune()
}

func (p *Parser) parseExpr(r0 rune) (Expr, error) {
	var id []rune
	for {
		r1, err := p.readRune()
		if err != nil {
			if err == io.EOF {
				err = UnexpectedEOF
			}
			return nil, err
		}

		if r1 == '}' {
			return Ident(string(id)), nil
		}

		if r1 == ':' {
			r2, err := p.readRune()
			if err != nil {
				return nil, err
			}
			switch r2 {
			case '-', '=':
				expr, err := p.parseFallback()
				return &Elvis{ident: Ident(string(id)), op: r2, fallback: expr}, err
			}
			p.unreadRune()
			return nil, UnexpectedRune
		}

		id = append(id, r1)
	}
}

func (p *Parser) parseFallback() (Expr, error) {
	var list List
	var err error
	for {
		var r0 rune
		r0, err = p.readRune()
		if err != nil {
			if err == io.EOF {
				err = UnexpectedEOF
			}
			break
		}

		if r0 == '}' {
			break
		}

		var expr Expr
		expr, err = p.parse(r0)
		if err != nil {
			break
		}
		list = append(list, expr)
	}
	return list, err
}

func (p *Parser) parseIdent(r0 rune) (Expr, error) {
	id := []rune{r0}
	var err error
	for {
		var r1 rune
		r1, err = p.readRune()
		if err != nil {
			break
		}

		if IsIdentRune(r1) {
			id = append(id, r1)
		} else {
			p.unreadRune()
			break
		}
	}
	return Ident(string(id)), err
}

func (p *Parser) parseEscape(r0 rune) (Expr, error) {
	r, err := p.readRune()
	if err != nil {
		r = r0
	}
	return Rune(r), err
}

func IsIdentHead(r rune) bool { return r == '_' || unicode.IsLetter(r) }

func IsIdentRune(r rune) bool { return IsIdentHead(r) || unicode.IsDigit(r) }
