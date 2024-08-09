package parse

import (
	"fmt"
	"iter"

	"alon.kr/x/view"
)

type unionParser[TokenT comparable, NodeT comparable] struct {
	name    string
	parsers []NodeParser[TokenT, NodeT]
}

func (p unionParser[TokenT, NodeT]) String() (s string) {
	for _, parser := range p.parsers[:len(p.parsers)-1] {
		s += parser.Name() + ", "
	}

	last := p.parsers[len(p.parsers)-1]
	return s + "or " + last.Name()
}

func (p unionParser[TokenT, NodeT]) Name() string {
	return p.name
}

func (p unionParser[TokenT, NodeT]) genError() error {
	return fmt.Errorf("expected %s", p.String())
}

func Or[TokenT comparable, NodeT comparable](
	name string,
	first, second NodeParser[TokenT, NodeT],
	additional ...NodeParser[TokenT, NodeT],
) NodeParser[TokenT, NodeT] {
	return unionParser[TokenT, NodeT]{
		name:    name,
		parsers: append([]NodeParser[TokenT, NodeT]{first, second}, additional...),
	}
}

func (p unionParser[TokenT, NodeT]) Parse(
	view *view.View[TokenT, uint],
) iter.Seq2[NodeT, error] {
	return func(yield func(NodeT, error) bool) {
		for _, parser := range p.parsers {
			for node, err := range parser.Parse(view) {
				if err == nil {
					if !yield(node, nil) {
						return
					}
				}
			}
		}

		var node NodeT
		yield(node, p.genError())
		return
	}
}
