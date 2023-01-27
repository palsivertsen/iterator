package iterator_test

import (
	"fmt"
	"testing"

	"iterator"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestIterator(t *testing.T) {
	t.Parallel()
	tests := []struct {
		items []string
	}{
		{},
		{
			items: []string{"a", "b", "c"},
		},
	}
	for testNum, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("next func: %d", testNum), func(t *testing.T) {
			t.Parallel()

			unit := iterator.Iterator[string]{
				NextFunc:    (&nextFuncHelper[string]{items: slices.Clone(tt.items)}).Items,
				ItemsBuffer: make([]string, 0, 1),
			}

			var items []string
			for unit.Next() {
				var item string
				unit.Item(&item)
				items = append(items, item)
			}

			assert.NoError(t, unit.Err())
			assert.Equal(t, tt.items, items)
		})
	}
	for testNum, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("no next func: %d", testNum), func(t *testing.T) {
			t.Parallel()

			unit := iterator.Iterator[string]{
				ItemsBuffer: slices.Clone(tt.items),
			}

			var items []string
			for unit.Next() {
				var item string
				unit.Item(&item)
				items = append(items, item)
			}

			assert.NoError(t, unit.Err())
			assert.Equal(t, tt.items, items)
		})
	}
}

type nextFuncHelper[T any] struct {
	items []T

	index int
}

func (nf *nextFuncHelper[T]) Items(items []T) (int, error) {
	if nf.index >= len(nf.items) {
		return 0, nil
	}

	i := copy(items, nf.items[nf.index:len(nf.items)])

	nf.index += i
	return i, nil
}
