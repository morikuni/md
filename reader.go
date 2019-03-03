package md

import (
	"bufio"
	"io"

	"golang.org/x/xerrors"
)

type LineReader struct {
	r    *bufio.Scanner
	peek string
	err  error
}

func NewLineReader(r io.Reader) *LineReader {
	return &LineReader{bufio.NewScanner(r), "", nil}
}

func (r *LineReader) PeekLine() (string, error) {
	if r.peek != "" || r.err != nil {
		return r.peek, r.err
	}
	r.peek, r.err = r.readLine()
	return r.peek, r.err
}

func (r *LineReader) Advance() {
	r.peek = ""
	r.err = nil
}

func (r *LineReader) readLine() (string, error) {
	if !r.r.Scan() {
		if err := r.r.Err(); err != nil {
			return "", xerrors.Errorf("read line: %w", err)
		}
		return "", xerrors.Errorf("read line: %w", io.EOF)
	}
	return r.r.Text(), nil
}
