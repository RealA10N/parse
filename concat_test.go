package parse_test

import (
	"iter"
	"testing"

	"alon.kr/x/parse"
	"alon.kr/x/view"
	"github.com/stretchr/testify/assert"
)

type Token uint8

type child = int
type father = []int

type positiveParser struct{}

func (positiveParser) String() string {
	return "+"
}

func (positiveParser) Parse(view *view.View[int, uint]) iter.Seq[child] {
	return func(yield func(child) bool) {
		n, err := view.At(0)
		if err != nil {
			return
		}
		if n > 0 {
			*view = view.Subview(1, view.Len())
			yield(n)
		}
	}
}

type negativeParser struct{}

func (negativeParser) String() string {
	return "-"
}

func (negativeParser) Parse(view *view.View[int, uint]) iter.Seq[child] {
	return func(yield func(child) bool) {
		n, err := view.At(0)
		if err != nil {
			return
		}
		if n < 0 {
			*view = view.Subview(1, view.Len())
			yield(n)
		}
	}
}

func TestConcatSimpleCase(t *testing.T) {
	numbersParser := parse.Concat(
		func(nodes []child) father { return father(nodes) },
		positiveParser{},
		negativeParser{},
	)

	t.Log(numbersParser.String())
	v := view.NewView[int, uint]([]int{1, -2})
	iter := numbersParser.Parse(&v)
	results := [][]int{}
	for nodes := range iter {
		results = append(results, nodes)
	}

	expected := [][]int{
		{1, -2},
	}

	assert.EqualValues(t, expected, results)
}
