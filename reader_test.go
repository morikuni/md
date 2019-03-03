package md

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLineReader(t *testing.T) {
	r := NewLineReader(strings.NewReader(`aaa

bbb`))

	l, err := r.PeekLine()
	require.NoError(t, err)
	require.Equal(t, "aaa", l)

	l, err = r.PeekLine()
	require.NoError(t, err)
	require.Equal(t, "aaa", l)

	r.Advance()

	l, err = r.PeekLine()
	require.NoError(t, err)
	require.Equal(t, "", l)

	r.Advance()

	l, err = r.PeekLine()
	require.NoError(t, err)
	require.Equal(t, "bbb", l)
}
