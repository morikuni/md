package md

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	cases := map[string]struct {
		file string

		wantElements []Element
		wantErr      bool
	}{
		"success": {
			"1.md",

			[]Element{
				&Header{Level: 1, Text: "aaa\nbbb"},
				&Header{Level: 2, Text: "ccc"},
				&Paragraph{Elements: []ParagraphElement{
					Text("paragraph1"), Code("code"),
					Text("\nparagraph2\n- paragraph5"),
				}},
				&List{
					Elements: []*ListElement{
						{1, "l1"},
						{2, "l2"},
						{1, "l3"},
					},
				},
				&CodeBlock{
					Language: "go",
					Code: `func main() {
	fmt.Println()
}`,
				},
				&Paragraph{Elements: []ParagraphElement{Text("paragraph3")}},
				&Quote{Text: "quote1\nquote2"},
				&Paragraph{Elements: []ParagraphElement{Text("paragraph4")}},
			},
			false,
		},
		"code element is not closed": {
			"2.md",

			nil,
			true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", tc.file))
			require.NoError(t, err)
			defer f.Close()

			es, err := Parse(f)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantElements, es)
		})
	}
}
