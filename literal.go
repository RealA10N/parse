package parse

import (
	"fmt"
	"iter"

	"alon.kr/x/view"
	"golang.org/x/exp/constraints"
)

type LiteralPropagator[TokenT any, NodeT any] func(TokenT) NodeT

type literalRule[
	TokenT comparable,
	NodeT any,
	OffsetT constraints.Unsigned,
] struct {
	propagator LiteralPropagator[TokenT, NodeT]
	literal    TokenT
}

func Literal[
	TokenT comparable,
	NodeT any,
	OffsetT constraints.Unsigned,
](
	propagator LiteralPropagator[TokenT, NodeT],
	token TokenT,
) Rule[TokenT, NodeT, OffsetT] {
	return literalRule[TokenT, NodeT, OffsetT]{
		propagator: propagator,
		literal:    token,
	}
}

func (p literalRule[TokenT, NodeT, OffsetT]) String() string {
	return fmt.Sprintf("%v", p.literal)
}

func (p literalRule[TokenT, NodeT, OffsetT]) Parse(
	v *view.View[TokenT, OffsetT],
) iter.Seq[NodeT] {
	return func(yield func(NodeT) bool) {
		token, err := v.At(0)
		if err == nil {
			yield(p.propagator(token))
		}
	}
}
