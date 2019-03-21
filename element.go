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
	Elements []ParagraphElement
}

func (*Paragraph) element() {}

type CodeBlock struct {
	Language string
	Code     string
}

func (*CodeBlock) element() {}

type ParagraphElement interface {
	textElement()
}

var _ = []ParagraphElement{
	Text(""),
	Code(""),
}

type Text string

func (Text) textElement() {}

type Code string

func (Code) textElement() {}
