package rules

import (
	"iter"

	"alon.kr/x/view"
	"golang.org/x/exp/constraints"
)

type UnionNodePropagator[ChildNodeT any, NodeT any] func(node ChildNodeT) NodeT

type unionRule[
	TokenT comparable,
	NodeT any,
	ChildNodeT any,
	OffsetT constraints.Unsigned,
] struct {
	propagator UnionNodePropagator[ChildNodeT, NodeT]
	rules      []Rule[TokenT, ChildNodeT, OffsetT]
}

func Union[
	TokenT comparable,
	NodeT comparable,
	ChildNodeT comparable,
	OffsetT constraints.Unsigned,
](
	propagator UnionNodePropagator[ChildNodeT, NodeT],
	first, second Rule[TokenT, ChildNodeT, OffsetT],
	additional ...Rule[TokenT, ChildNodeT, OffsetT],
) Rule[TokenT, NodeT, OffsetT] {
	return unionRule[TokenT, NodeT, ChildNodeT, OffsetT]{
		propagator: propagator,
		rules:      append([]Rule[TokenT, ChildNodeT, OffsetT]{first, second}, additional...),
	}
}

func (r unionRule[TokenT, NodeT, ChildNodeT, OffsetT]) String() string {
	s := "(" + r.rules[0].String()
	for _, rule := range r.rules[1:] {
		s += " | " + rule.String()
	}
	return s + ")"
}

func (r unionRule[TokenT, NodeT, ChildNodeT, OffsetT]) Parse(
	v *view.View[TokenT, OffsetT],
) iter.Seq[NodeT] {
	return func(yield func(NodeT) bool) {
		for _, rule := range r.rules {
			for node := range rule.Parse(v) {
				if !yield(r.propagator(node)) {
					return
				}
			}
		}
	}
}
