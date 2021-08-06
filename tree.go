package webkit

import (
	"net/http"
	"sort"
)

type RouteType uint8

const (
	RouteDynamic RouteType = 1 << iota
	RouteWild
	RouteStatic
)

type nodes []*node

func (a nodes) Len() int           { return len(a) }
func (a nodes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a nodes) Less(i, j int) bool { return a[i].typ > a[j].typ }

type node struct {
	typ      RouteType
	segment  string
	handler  http.Handler
	children nodes
}

func (n *node) insert(segments []string, handler http.Handler, height int) {
	if height == len(segments) {
		n.handler = handler
		return
	}

	segment := segments[height]
	child := n.find(segment)

	if child == nil {
		child = &node{
			typ:     weighten(segment),
			segment: segment,
		}
		// lazy init, cuz only internal nodes need allocation for children
		if n.children == nil {
			n.children = make(nodes, 0)
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

func (n *node) lookup(segments []string, height int, params *Params, target *node) {
	if height == len(segments) ||
		(len(n.segment) != 0 && n.segment[0] == '*') {
		*target = *n
		return
	}

	var route RouteType
	segment := segments[height]

	for _, child := range n.children {
		switch {
		case child.segment == segment:
			route = RouteStatic
		case child.segment[0] == '*':
			route = RouteWild
		case child.segment[0] == ':':
			route = RouteDynamic
		}

		if route == RouteDynamic {
			if *params == nil {
				*params = make(Params)
			}
			(*params)[child.segment[1:]] = segment
		}

		if route != 0 {
			if child.lookup(segments, height+1, params, target); target != nil {
				return
			}
		}

		route = 0
	}
}

func weighten(s string) RouteType {
	switch s[0] {
	case ':':
		return RouteDynamic
	case '*':
		return RouteWild
	default:
		return RouteStatic
	}
}
