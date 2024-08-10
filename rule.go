package parse

import (
	"iter"

	"alon.kr/x/view"
	"golang.org/x/exp/constraints"
)

type Rule[
	TokenT comparable,
	NodeT any,
	OffsetT constraints.Unsigned,
] interface {
	String() string
	Parse(view *view.View[TokenT, OffsetT]) iter.Seq[NodeT]
}
