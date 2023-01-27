package iterator_test

import (
	"fmt"
	"strconv"
	"testing"

	"iterator"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestIterator_ZeroAllocations(t *testing.T) {
	res := testing.Benchmark(func(b *testing.B) {
		nfh := nextFuncHelper[string]{items: make([]string, b.N)}
		unit := iterator.Iterator[string]{
			NextFunc:    nfh.Items,
			ItemsBuffer: make([]string, 1),
		}
		b.ResetTimer()

		var t string
		for i := 0; i < b.N; i++ {
			unit.Next()
			unit.Item(&t)
		}
		b.StopTimer()
		require.NoError(b, unit.Err())
	})
	assert.Equal(t, int64(0), res.AllocsPerOp())
}

func BenchmarkIterator(b *testing.B) {
	benchmarks := []struct {
		bufferSize int
	}{
		{0},
		{1 << 0},
		{1 << 1},
		{1 << 2},
		{1 << 3},
		{1 << 4},
		{1 << 5},
		{1 << 6},
		{1 << 7},
		{1 << 8},
		{1 << 9},
		{1 << 10},
		{1 << 11},
		{1 << 12},
	}
	for _, bm := range benchmarks {
		bm := bm
		b.Run(strconv.Itoa(bm.bufferSize), func(b *testing.B) {
			nfh := nextFuncHelper[struct{}]{items: make([]struct{}, b.N)}
			unit := iterator.Iterator[struct{}]{
				NextFunc:    nfh.Items,
				ItemsBuffer: make([]struct{}, bm.bufferSize),
			}
			b.ResetTimer()

			var t struct{}
			for i := 0; i < b.N; i++ {
				unit.Next()
				unit.Item(&t)
			}
			b.StopTimer()
			require.NoError(b, unit.Err())
		})
	}
}

// helpers

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
