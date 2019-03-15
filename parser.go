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
			result, err = flushParagraph(result, paragraph)
			if err != nil {
				return nil, err
			}
			paragraph = nil
		}
	}

	var err error
	result, err = flushParagraph(result, paragraph)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func flushParagraph(r []Element, paragraph []string) ([]Element, error) {
	if paragraph == nil {
		return r, nil
	}
	p, err := convertParagraph(strings.Join(paragraph, "\n"))
	if err != nil {
		return nil, err
	}
	return append(r, p), nil
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

type pState int

const (
	text pState = iota
	code
)

func convertParagraph(paragraph string) (*Paragraph, error) {
	runes := []rune(paragraph)

	var (
		from     int
		state    pState
		elements []ParagraphElement
	)

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch r {
		case '`':
			switch state {
			case text:
				elements = append(elements, Text(runes[from:i]))
				state = code
			case code:
				elements = append(elements, Code(runes[from:i]))
				state = text
			}
			from = i + 1
		default:
		}
	}
	if from < len(runes) {
		switch state {
		case text:
			elements = append(elements, Text(runes[from:]))
		case code:
			return nil, xerrors.Errorf("code element should be closed around: %q", string(runes[from:]))
		}
	}
	return &Paragraph{elements}, nil
}
