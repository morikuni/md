package md

import (
	"io"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/xerrors"
)

func Parse(r io.Reader) ([]Element, error) {
	return parse(NewLineReader(r))
}

func parse(r *lineReader) ([]Element, error) {
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
		case isEmpty(line):
			// ignore
			r.Advance()
		case line[0] == '#':
			h, err := readHeader(r)
			if err != nil {
				return nil, err
			}
			result = append(result, h)
		case strings.HasPrefix(line, "```"):
			cb, err := readCodeBlock(r)
			if err != nil {
				return nil, err
			}
			result = append(result, cb)
		case strings.HasPrefix(strings.TrimLeft(line, " \t"), "-"):
			li, err := readList(r)
			if err != nil {
				return nil, err
			}
			result = append(result, li)
		case strings.HasPrefix(line, ">"):
			cb, err := readQuote(r)
			if err != nil {
				return nil, err
			}
			result = append(result, cb)
		default:
			p, err := readParagraph(r)
			if err != nil {
				return nil, err
			}
			result = append(result, p)
		}
	}

	return result, nil
}

func readHeader(r *lineReader) (*Header, error) {
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

func readCodeBlock(r *lineReader) (*CodeBlock, error) {
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

func readList(r *lineReader) (*List, error) {
	var elements []ListElement
	for {
		line, err := r.PeekLine()
		if xerrors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(strings.TrimLeft(line, " \t"), "-") {
			break
		}
		r.Advance()
		level := countLeft(line, ' ')/2 + countLeft(line, '\t') + 1
		text := strings.TrimLeft(line, " \t-")
		elements = append(elements, ListElement{
			level,
			text,
		})
	}
	return &List{elements}, nil
}

func readParagraph(r *lineReader) (*Paragraph, error) {
	var lines []string
	for {
		line, err := r.PeekLine()
		if xerrors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if isEmpty(line) {
			break
		}
		r.Advance()
		lines = append(lines, line)
	}

	tes, err := parseTextBlock(strings.Join(lines, "\n"))
	if err != nil {
		return nil, err
	}
	return &Paragraph{tes}, nil
}

var (
	reCode = regexp.MustCompile("\\A`(.*?)[^\\\\]`")
	reLink = regexp.MustCompile(`\A\[(.+?)\]\((.+?)\)`)
)

func parseTextBlock(paragraph string) (*TextBlock, error) {
	data := []byte(paragraph)

	var (
		elements []TextElement
		last     rune
	)

	var text []byte
	flush := func() {
		if len(text) == 0 {
			return
		}
		elements = append(elements, Text(text))
		text = text[:0]
	}
	dropRight := func() {
		if len(text) == 0 {
			return
		}
		text = text[:len(text)-1]
	}
	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		switch {
		case last == '\\' && strings.ContainsAny(string(r), "`\\["):
			dropRight()
			fallthrough
		default:
			text = append(text, data[:size]...)
		case r == '`':
			idx := reCode.FindSubmatchIndex(data)
			if len(idx) == 0 {
				break
			}
			flush()
			from, to := idx[2], idx[3]
			size = idx[1]
			// +1 for a character matched to [^\\\\]
			elements = append(elements, Code(data[from:to+1]))
		case r == '[':
			idx := reLink.FindSubmatchIndex(data)
			if len(idx) == 0 {
				break
			}
			flush()
			text := data[idx[2]:idx[3]]
			ref := data[idx[4]:idx[5]]
			size = idx[1]
			elements = append(elements, Link{string(text), string(ref)})
		}
		last = r
		data = data[size:]
	}
	flush()

	return &TextBlock{elements}, nil
}

func readQuote(r *lineReader) (*Quote, error) {
	var lines []string
	for {
		line, err := r.PeekLine()
		if xerrors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(line, ">") {
			break
		}
		r.Advance()
		text := strings.TrimLeft(line, "> ")
		r.Advance()
		lines = append(lines, text)
	}
	return &Quote{strings.Join(lines, "\n")}, nil
}
