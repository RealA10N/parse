package parse

import (
	"iter"

	"alon.kr/x/view"
)

type expect[TokenT comparable, ErrorT any] struct {
	tokens       view.View[TokenT, uint]
	errorFactory ErrorFactory[TokenT, ErrorT]
}

func (node *expect[TokenT, ErrorT]) Parse(
	v *view.View[TokenT, uint],
) iter.Seq2[view.UnmanagedView[TokenT, uint], ErrorT] {
	return func(yield func(view.UnmanagedView[TokenT, uint], ErrorT) bool) {
		if v.Len() < node.tokens.Len() {
			expected := []TokenT{node.tokens.AtUnsafe(v.Len() - 1)}
			var actual TokenT
			err := node.errorFactory(expected, actual)
			var u view.UnmanagedView[TokenT, uint]
			yield(u, err)
			return
		}

		longestPrefix := v.LongestCommonPrefix(node.tokens).Len()
		if longestPrefix < node.tokens.Len() {
			expected := []TokenT{node.tokens.AtUnsafe(longestPrefix)}
			actual := v.AtUnsafe(longestPrefix)
			err := node.errorFactory(expected, actual)
			var u view.UnmanagedView[TokenT, uint]
			yield(u, err)
			return
		}

		prefix, rest := v.Partition(node.tokens.Len())
		*v = rest
		unmanaged, _ := prefix.Detach()
		var err ErrorT
		yield(unmanaged, err)
	}
}

func Expect[TokenT comparable, ErrorT any](
	errorFactory ErrorFactory[TokenT, ErrorT],
	first TokenT,
	additional ...TokenT,
) Node[TokenT, ErrorT] {
	tokens := append([]TokenT{first}, additional...)
	return &expect[TokenT, ErrorT]{
		tokens:       view.NewView[TokenT, uint](tokens),
		errorFactory: errorFactory,
	}
}
