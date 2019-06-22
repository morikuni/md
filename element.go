package md

type Element interface {
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

func (*Header) element() {}

type CodeBlock struct {
	Language string
	Code     string
}

func (*CodeBlock) element() {}

type List struct {
	Elements []*ListElement
}

func (*List) element() {}

type ListElement struct {
	Level int
	Text  string
}

type Paragraph struct {
	Elements []ParagraphElement
}

func (*Paragraph) element() {}

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

type Quote struct {
	Text string
}

func (*Quote) element() {}
