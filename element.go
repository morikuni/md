package md

type Element interface {
	element()
}

var _ = []Element{
	(*Header)(nil),
	(*Paragraph)(nil),
	(*CodeBlock)(nil),
}

type Header struct {
	Level int
	Text  string
}

func (*Header) element() {}

type Paragraph struct {
	Elements []TextElement
}

func (*Paragraph) element() {}

type CodeBlock struct {
	Code     string
	Language string
}

func (*CodeBlock) element() {}

type TextElement interface {
	textElement()
}

var _ = []TextElement{
	Text(""),
	Code(""),
}

type Text string

func (Text) textElement() {}

type Code string

func (Code) textElement() {}
