package emitter

import (
	"strings"
	"sync"
)

// route is a value associated with a subscription.
type route struct {
	Topic  string
	Action MessageHandler
}

func newRoute(topic string, handler MessageHandler) ([]string, route) {
	query := strings.FieldsFunc(topic, func(c rune) bool {
		return c == '/'
	})

	return query, route{
		Topic:  topic,
		Action: handler,
	}
}

// ------------------------------------------------------------------------------------

type node struct {
	word     string
	routes   map[string]route
	parent   *node
	children map[string]*node
}

func (n *node) orphan() {
	if n.parent == nil {
		return
	}

	delete(n.parent.children, n.word)
	if len(n.parent.routes) == 0 && len(n.parent.children) == 0 {
		n.parent.orphan()
	}
}

// trie represents an efficient collection of subscriptions with lookup capability.
type trie struct {
	sync.RWMutex
	root *node // The root node of the tree.
	lookup func(query []string, result *[]MessageHandler, node *node)
}

// newTrie creates a new trie without a lookup function.
func newTrie() *trie {
	return &trie{
		root: &node{
			children: make(map[string]*node),
		},
	}
}

// NewTrie creates a new subscriptions matcher using standard emitter strategy.
func NewTrie() *trie {
	t := newTrie()
	t.lookup = t.lookupEmitter
	return t
}

// NewTrieMQTT creates a new subscriptions matcher using standard MQTT strategy.
func NewTrieMQTT() *trie {
	t := newTrie()
	t.lookup = t.lookupMqtt
	return t
}

// AddHandler adds a message handler to a topic.
func (t *trie) AddHandler(topic string, handler MessageHandler) error {
	query, rt := newRoute(topic, handler)

	t.Lock()
	curr := t.root
	for _, word := range query {
		child, ok := curr.children[word]
		if !ok {
			child = &node{
				word:     word,
				parent:   curr,
				routes:   make(map[string]route),
				children: make(map[string]*node),
			}
			curr.children[word] = child
		}
		curr = child
	}

	// Add the handler
	curr.routes[rt.Topic] = rt
	t.Unlock()
	return nil
}

// RemoveHandler removes a message handler from a topic.
func (t *trie) RemoveHandler(topic string) {
	query, _ := newRoute(topic, nil)

	t.Lock()
	curr := t.root
	for _, word := range query {
		child, ok := curr.children[word]
		if !ok {
			// Subscription doesn't exist.
			t.Unlock()
			return
		}
		curr = child
	}

	// Remove the route
	delete(curr.routes, topic)

	// Remove orphans
	if len(curr.routes) == 0 && len(curr.children) == 0 {
		curr.orphan()
	}
	t.Unlock()
}

// Lookup returns the handlers for the given topic.
func (t *trie) Lookup(topic string) []MessageHandler {
	query, _ := newRoute(topic, nil)
	var result []MessageHandler

	t.RLock()
	t.lookup(query, &result, t.root)
	t.RUnlock()
	return result
}

func (t *trie) lookupEmitter(query []string, result *[]MessageHandler, node *node) {

	// Add routes from the current branch
	for _, route := range node.routes {
		*result = append(*result, route.Action)
	}

	// If we're not yet done, continue
	if len(query) > 0 {

		// Go through the exact match branch
		if n, ok := node.children[query[0]]; ok {
			t.lookupEmitter(query[1:], result, n)
		}

		// Go through wildcard match branch
		if n, ok := node.children["+"]; ok {
			t.lookupEmitter(query[1:], result, n)
		}
	}
}

func (t *trie) lookupMqtt(query []string, result *[]MessageHandler, node *node) {
	if len(query) == 0 {
		// Add routes from the current branch
		for _, route := range node.routes {
			*result = append(*result, route.Action)
		}
	}

	// If we're not yet done, continue
	if len(query) > 0 {
		// Go through the exact match branch
		if n, ok := node.children[query[0]]; ok {
			t.lookupMqtt(query[1:], result, n)
		}

		// Go through wildcard match branch
		if n, ok := node.children["+"]; ok {
			t.lookupMqtt(query[1:], result, n)
		}

		// Go through wildcard match branch
		if n, ok := node.children["#"]; ok {
			for _, route := range n.routes {
				*result = append(*result, route.Action)
			}
		}
	}
}
