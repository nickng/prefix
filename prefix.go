package prefix

import (
	"fmt"
	"io"
)

type reader struct {
	rd  io.Reader
	buf []byte
}

func NewReader(rd io.Reader, rw Rewriter) (*reader, error) {
	var buf []byte
	for {
		p := make([]byte, 1024)
		n, err := rd.Read(p)
		if err != nil {
			return &reader{rd: rd, buf: p[:n]}, err
		}
		rewritten, done := rw.Rewrite(p, n)
		fmt.Println("Rewritten", string(p), " to", string(rewritten))
		buf = append(buf, rewritten...)
		if done {
			break
		}
	}
	return &reader{rd: rd, buf: buf}, nil
}

func (r *reader) Read(p []byte) (int, error) {
	if len(r.buf) > 0 {
		ncopied := copy(p, r.buf)
		r.buf = r.buf[ncopied:]

		// p is full, we can't copy all from buf to p
		if len(r.buf) > 0 {
			return ncopied, nil
		}

		// p has remaining space, so we fill it up with more Read
		readMore := make([]byte, len(p)-ncopied)
		nmore, err := r.rd.Read(readMore)
		if err != nil {
			return ncopied + nmore, err
		}
		copy(p[ncopied:], readMore)
		return ncopied + nmore, nil
	}

	return r.rd.Read(p)
}

// Rewriter is an interface for prefix rewriter.
//
// A prefix rewriter will keep calling Rewrite until it returns true
// and replace the bytes read with the rewritten bytes returned by it.
type Rewriter interface {
	// Rewrite should return true when the rewrite is complete
	// and we don't need to mutate the stream any more.
	Rewrite([]byte, int) ([]byte, bool)
}
