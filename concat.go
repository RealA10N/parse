package parse

import (
	"iter"

	"alon.kr/x/view"
	"golang.org/x/exp/constraints"
)

type ConcatNodePropagation[ChildNodeT any, NodeT any] func(nodes []ChildNodeT) NodeT

type concatParser[
	TokenT comparable,
	NodeT any,
	ChildNodeT any,
	OffsetT constraints.Unsigned,
] struct {
	propagator ConcatNodePropagation[ChildNodeT, NodeT]
	parsers    []NodeParser[TokenT, ChildNodeT, OffsetT]
}

func Concat[
	TokenT comparable,
	NodeT any,
	ChildNodeT any,
	OffsetT constraints.Unsigned,
](
	propagator ConcatNodePropagation[ChildNodeT, NodeT],
	first, second NodeParser[TokenT, ChildNodeT, OffsetT],
	additional ...NodeParser[TokenT, ChildNodeT, OffsetT],
) NodeParser[TokenT, NodeT, OffsetT] {
	return concatParser[TokenT, NodeT, ChildNodeT, OffsetT]{
		propagator: propagator,
		parsers: append(
			[]NodeParser[TokenT, ChildNodeT, OffsetT]{first, second},
			additional...,
		),
	}
}

func (p concatParser[TokenT, NodeT, ChildNodeT, OffsetT]) String() string {
	s := "("
	for _, parser := range p.parsers[:len(p.parsers)-1] {
		s += parser.String() + " "
	}
	return s + p.parsers[len(p.parsers)-1].String() + ")"
}

func (p concatParser[TokenT, NodeT, ChildNodeT, OffsetT]) parseSuffix(
	v *view.View[TokenT, OffsetT],
	collectedPrefix []ChildNodeT,
	k OffsetT,
) iter.Seq2[[]ChildNodeT, Error[OffsetT]] {
	return func(yield func([]ChildNodeT, Error[OffsetT]) bool) {
		n := OffsetT(len(p.parsers))

		// recursion base:
		// if already collected the whole prefix, yield and return
		if k == n {
			yield(collectedPrefix, nil)
			return
		}

		// save the current view state on the stack (push-down automaton!)
		// and parse the next node. Use the bookmark to restore the view state
		// in the case consecutive parsers fail and we need to backtrack.
		bookmark := *v
		curParser := p.parsers[k]

		for curNode, err := range curParser.Parse(v) {
			// if the current parser fails, backtrack.
			if err != nil {
				break
			}

			// if current parser parsed a node successfully, append it to
			// the existing suffix and recurse by yielding all values from
			// the recursive call.
			collectedPrefix[k] = curNode
			for node, err := range p.parseSuffix(v, collectedPrefix, k+1) {
				if !yield(node, err) || err != nil {
					return
				}
			}
		}

		*v = bookmark
	}
}

func (p concatParser[TokenT, NodeT, ChildNodeT, OffsetT]) Parse(
	v *view.View[TokenT, OffsetT],
) iter.Seq2[NodeT, Error[OffsetT]] {
	return func(yield func(NodeT, Error[OffsetT]) bool) {
		n := len(p.parsers)
		nodes := make([]ChildNodeT, n)
		for childNodes, err := range p.parseSuffix(v, nodes, 0) {
			if !yield(p.propagator(childNodes), err) || err != nil {
				return
			}
		}
	}
}
