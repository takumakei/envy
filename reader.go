package main

import "io"

type Reader struct {
	list [][]byte
}

func NewReader(data ...[]byte) *Reader {
	return &Reader{list: data}
}

func (r *Reader) Read(p []byte) (int, error) {
	if len(r.list) == 0 {
		return 0, io.EOF
	}
	s := r.list[0]
	slen := len(s)
	plen := len(p)
	copy(p, s)
	if slen <= plen {
		r.list = r.list[1:]
		return slen, nil
	} else {
		r.list[0] = s[plen:]
		return plen, nil
	}
}
