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
		result []Element
	)
	for {
		line, err := r.PeekLine()
		if err != nil {
			if xerrors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		switch {
		case len(line) == 0:
			// ignore
			r.Advance()
		case line[0] == '#': // header
			h, err := readHeader(r)
			if err != nil {
				return nil, err
			}
			result = append(result, h)
		case strings.HasPrefix(line, "```"): // code block
			cb, err := readCodeBlock(r)
			if err != nil {
				return nil, err
			}
			result = append(result, cb)
		default: // paragraph
			p, err := readParagraph(r)
			if err != nil {
				return nil, err
			}
			result = append(result, p)
		}
	}

	return result, nil
}

func readCodeBlock(r *LineReader) (*CodeBlock, error) {
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
		lang,
		strings.Join(codes, "\n"),
	}, nil
}

func readHeader(r *LineReader) (*Header, error) {
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

func readParagraph(r *LineReader) (*Paragraph, error) {
	var lines []string
	for {
		line, err := r.PeekLine()
		if xerrors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if line == "" {
			break
		}
		r.Advance()
		lines = append(lines, line)
	}
	return convertParagraph(strings.Join(lines, "\n"))
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
