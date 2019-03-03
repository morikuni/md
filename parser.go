package md

import (
	"io"
	"strings"

	"golang.org/x/xerrors"
)

func Parse(r io.Reader) ([]Element, error) {
	return parse(NewLineReader(r))
}

func parse(r *LineReader) ([]Element, error) {
	var (
		result    []Element
		paragraph []string
	)
	for {
		line, err := r.PeekLine()
		if err != nil {
			if xerrors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		isParagraph := false
		switch {
		case len(line) == 0:
			// ignore
			r.Advance()
		case line[0] == '#': // header
			cb, err := readHeader(r)
			if err != nil {
				return nil, err
			}
			result = append(result, cb)
		case strings.HasPrefix(line, "```"): // code block
			cb, err := readCodeBlock(r)
			if err != nil {
				return nil, err
			}
			result = append(result, cb)
		default: // paragraph
			isParagraph = true
			paragraph = append(paragraph, line)
			r.Advance()
		}

		if !isParagraph {
			result = flushParagraph(result, paragraph)
			paragraph = nil
		}
	}

	result = flushParagraph(result, paragraph)

	return result, nil
}

func flushParagraph(r []Element, paragraph []string) []Element {
	if paragraph == nil {
		return r
	}
	return append(r, &Paragraph{
		[]TextElement{Text(strings.Join(paragraph, "\n"))},
	})
}

func readCodeBlock(r *LineReader) (cb *CodeBlock, err error) {
	lang := strings.TrimLeft(r.MustPeekLine(), "`")
	r.Advance()

	var codes []string
	for {
		line, err := r.PeekLine()
		if err != nil {
			return nil, err
		}
		r.Advance()
		if line == "```" {
			break
		}
		codes = append(codes, line)
	}

	return &CodeBlock{
		strings.Join(codes, "\n"),
		lang,
	}, nil
}

func readHeader(r *LineReader) (cb *Header, err error) {
	level := countLeft(r.MustPeekLine(), '#')

	var headers []string
	for {
		line, err := r.PeekLine()
		if err != nil {
			return nil, err
		}
		l := countLeft(line, '#')
		if level != l {
			break
		}
		r.Advance()
		header := strings.TrimLeft(line, "# ")
		headers = append(headers, header)
	}

	return &Header{
		level,
		strings.Join(headers, "\n"),
	}, nil
}

func countLeft(s string, r rune) int {
	var count int
	for _, sr := range s {
		if sr != r {
			return count
		}
		count++
	}
	return count
}
