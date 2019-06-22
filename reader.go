package md

import (
	"bufio"
	"io"

	"golang.org/x/xerrors"
)

type lineReader struct {
	r    *bufio.Scanner
	peek string
	err  error
}

func NewLineReader(r io.Reader) *lineReader {
	return &lineReader{bufio.NewScanner(r), "", nil}
}

func (r *lineReader) PeekLine() (string, error) {
	if r.peek != "" || r.err != nil {
		return r.peek, r.err
	}
	r.peek, r.err = r.readLine()
	return r.peek, r.err
}

func (r *lineReader) MustPeekLine() string {
	l, err := r.PeekLine()
	if err != nil {
		panic(err)
	}
	return l
}

func (r *lineReader) Advance() {
	r.peek = ""
	r.err = nil
}

func (r *lineReader) readLine() (string, error) {
	if !r.r.Scan() {
		if err := r.r.Err(); err != nil {
			return "", xerrors.Errorf("read line: %w", err)
		}
		return "", xerrors.Errorf("read line: %w", io.EOF)
	}
	return r.r.Text(), nil
}
