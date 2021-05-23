package webkit

import (
	"fmt"
	"net/http"
	"sort"
)

const (
	TypeDynamic uint8 = 1 << iota
	TypeWild
	TypeStatic
)

type sortedNodes []*node

func (a sortedNodes) Len() int           { return len(a) }
func (a sortedNodes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortedNodes) Less(i, j int) bool { return a[i].weight > a[j].weight }

type node struct {
	weight   uint8
	segment  string
	handler  http.Handler
	children sortedNodes
}

func (n *node) insert(segments []string, handler http.Handler, height int) {
	if len(segments) == height {
		n.handler = handler
		return
	}
	segment := segments[height]

	child := n.find(segment)
	if child == nil {
		child = &node{weight: weight(segment), segment: segment, children: nil, handler: nil}
		if n.children == nil {
			n.children = make(sortedNodes, 0)
		}
		n.children = append(n.children, child)
		sort.Sort(n.children)
	}
	child.insert(segments, handler, height+1)
}

func (n *node) find(segment string) *node {
	for _, child := range n.children {
		if child.segment == segment {
			return child
		}
	}
	return nil
}

func (n *node) lookup(segments []string, height int, params *Params) *node {
	if len(segments) == height || (len(n.segment) != 0 && n.segment[0] == '*') {
		return n
	}

	segment := segments[height]

	for _, child := range n.children {
		var flag uint8
		if segment == child.segment {
			flag |= TypeStatic
		}
		if child.segment[0] == '*' {
			flag |= TypeWild
		}
		if child.segment[0] == ':' {
			flag |= TypeDynamic
		}

		if (flag & TypeDynamic) != 0 {
			if *params == nil {
				*params = make(Params)
			}
			(*params)[child.segment[1:]] = segment
		}

		if flag != 0 {
			res := child.lookup(segments, height+1, params)
			if res != nil {
				return res
			}
		}

	}

	return nil
}

func (n *node) String() string {
	return fmt.Sprintf("%d %s %v %d", n.weight, n.segment, n.handler, len(n.children))
}

func weight(segment string) uint8 {
	switch segment[0] {
	case ':':
		return 1
	case '*':
		return 2
	default:
		return 3
	}
}
