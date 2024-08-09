package parse

import (
	"iter"

	"alon.kr/x/view"
)

type Node[TokenT comparable, ErrorT any] interface {
	Parse(view *view.View[TokenT, uint]) iter.Seq2[view.UnmanagedView[TokenT, uint], ErrorT]
}

type ErrorFactory[TokenT comparable, ErrorT any] func(expected []TokenT, actual TokenT) ErrorT
