package emitter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoute(t *testing.T) {
	query, _ := newRoute("a/", nil)
	assert.Len(t, query, 1)
}

func TestTrieMatch(t *testing.T) {
	m := NewTrie()
	testPopulateWithStrings(m, []string{
		"a/",
		"a/b/c/",
		"a/+/c/",
		"a/b/c/d/",
		"a/+/c/+/",
		"x/",
		"x/y/",
		"x/+/z",
	})

	// Tests to run
	tests := []struct {
		topic string
		n     int
	}{
		{topic: "a/", n: 1},
		{topic: "a/1/", n: 1},
		{topic: "a/2/", n: 1},
		{topic: "a/1/2/", n: 1},
		{topic: "a/1/2/3/", n: 1},
		{topic: "a/x/y/c/", n: 1},
		{topic: "a/x/c/", n: 2},
		{topic: "a/b/c/", n: 3},
		{topic: "a/b/c/d/", n: 5},
		{topic: "a/b/c/e/", n: 4},
		{topic: "x/y/c/e/", n: 2},
	}

	for _, tc := range tests {
		result := m.Lookup(tc.topic)
		assert.Equal(t, tc.n, len(result))
	}
}

// Populates the trie with a set of strings
func testPopulateWithStrings(m *trie, values []string) {
	for _, s := range values {
		m.AddHandler(s, func(_ *Client, _ Message) {
			println("dummy handler")
		})
	}
}
