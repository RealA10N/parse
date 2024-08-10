package rules

import (
	"iter"

	"alon.kr/x/view"
	"golang.org/x/exp/constraints"
)

type ConcatNodePropagator[ChildNodeT any, NodeT any] func(nodes []ChildNodeT) NodeT

type concatRule[
	TokenT comparable,
	NodeT any,
	ChildNodeT any,
	OffsetT constraints.Unsigned,
] struct {
	propagator ConcatNodePropagator[ChildNodeT, NodeT]
	rules      []Rule[TokenT, ChildNodeT, OffsetT]
}

func Concat[
	TokenT comparable,
	NodeT any,
	ChildNodeT any,
	OffsetT constraints.Unsigned,
](
	propagator ConcatNodePropagator[ChildNodeT, NodeT],
	first, second Rule[TokenT, ChildNodeT, OffsetT],
	additional ...Rule[TokenT, ChildNodeT, OffsetT],
) Rule[TokenT, NodeT, OffsetT] {
	return concatRule[TokenT, NodeT, ChildNodeT, OffsetT]{
		propagator: propagator,
		rules: append(
			[]Rule[TokenT, ChildNodeT, OffsetT]{first, second},
			additional...,
		),
	}
}

func (r concatRule[TokenT, NodeT, ChildNodeT, OffsetT]) String() string {
	s := "(" + r.rules[0].String()
	for _, parser := range r.rules[1:] {
		s += " " + parser.String()
	}
	return s + ")"
}

func (r concatRule[TokenT, NodeT, ChildNodeT, OffsetT]) parseSuffix(
	v *view.View[TokenT, OffsetT],
	collectedPrefix []ChildNodeT,
	k OffsetT,
) iter.Seq[[]ChildNodeT] {
	return func(yield func([]ChildNodeT) bool) {
		n := OffsetT(len(r.rules))

		// recursion base:
		// if already collected the whole prefix, yield and return
		if k == n {
			yield(collectedPrefix)
			return
		}

		// save the current view state on the stack (push-down automaton!)
		// and parse the next node. Use the bookmark to restore the view state
		// in the case consecutive parsers fail and we need to backtrack.
		bookmark := *v
		curParser := r.rules[k]

		for curNode := range curParser.Parse(v) {
			// if current parser parsed a node successfully, append it to
			// the existing suffix and recurse by yielding all values from
			// the recursive call.
			collectedPrefix[k] = curNode
			for node := range r.parseSuffix(v, collectedPrefix, k+1) {
				if !yield(node) {
					return
				}
			}
		}

		*v = bookmark
	}
}

func (r concatRule[TokenT, NodeT, ChildNodeT, OffsetT]) Parse(
	v *view.View[TokenT, OffsetT],
) iter.Seq[NodeT] {
	return func(yield func(NodeT) bool) {
		n := len(r.rules)
		nodes := make([]ChildNodeT, n)
		for childNodes := range r.parseSuffix(v, nodes, 0) {
			if !yield(r.propagator(childNodes)) {
				return
			}
		}
	}
}
