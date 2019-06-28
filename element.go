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
	TextBlock *TextBlock
}

func (*Paragraph) element() {}

type TextBlock struct {
	Elements []TextElement
}

type TextElement interface {
	textElement()
}

var _ = []TextElement{
	Text(""),
	Code(""),
	Link{},
}

type Text string

func (Text) textElement() {}

type Code string

func (Code) textElement() {}

type Link struct {
	Text      string
	Reference string
}

func (Link) textElement() {}

type Quote struct {
	Text string
}

func (*Quote) element() {}
