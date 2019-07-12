package md

import (
	"fmt"
	"strings"
)

type Element interface {
	fmt.Stringer
	element()
}

var _ = []Element{
	(*Header)(nil),
	(*Paragraph)(nil),
	(*CodeBlock)(nil),
	(*List)(nil),
	(*Quote)(nil),
}

type Header struct {
	Level int
	Text  string
}

func (h *Header) String() string {
	return h.Text
}

func (*Header) element() {}

type CodeBlock struct {
	Language string
	Code     string
}

func (cb *CodeBlock) String() string {
	return cb.Code
}

func (*CodeBlock) element() {}

type List struct {
	Elements []ListElement
}

func (h *List) String() string {
	var b strings.Builder
	first := true
	for _, e := range h.Elements {
		if !first {
			b.WriteRune(' ')
		}
		b.WriteString(e.String())
	}
	return b.String()
}

func (*List) element() {}

type ListElement struct {
	Level int
	Text  string
}

func (le ListElement) String() string {
	return le.Text
}

type Paragraph struct {
	TextBlock *TextBlock
}

func (p *Paragraph) String() string {
	return p.TextBlock.String()
}

func (*Paragraph) element() {}

type TextBlock struct {
	Elements []TextElement
}

func (tb *TextBlock) String() string {
	var b strings.Builder
	for _, te := range tb.Elements {
		b.WriteString(te.String())
	}
	return b.String()
}

type TextElement interface {
	fmt.Stringer
	textElement()
}

var _ = []TextElement{
	Text(""),
	Code(""),
	Link{},
}

type Text string

func (t Text) String() string {
	return string(t)
}

func (Text) textElement() {}

type Code string

func (c Code) String() string {
	return string(c)
}

func (Code) textElement() {}

type Link struct {
	Text      string
	Reference string
}

func (l Link) String() string {
	return string(l.Text)
}

func (Link) textElement() {}

type Quote struct {
	Text string
}

func (q *Quote) String() string {
	return q.Text
}

func (*Quote) element() {}
