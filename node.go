package parse

import (
	"iter"

	"alon.kr/x/view"
)

type NodeParser[TokenT comparable, NodeT any] interface {
	String() string
	Name() string
	Parse(view *view.View[TokenT, uint]) iter.Seq2[NodeT, error]
}
