// Package iterator provides a simple way of working with pagination APIs.
package iterator

import (
	"errors"
	"fmt"
)

const (
	errEOI = internalError("end of items")
)

// NextPageFunc copies at most len(t) items into t. A return value of 0
// indicates that there are no more items.
type NextPageFunc[T any] func(t []T) (int, error)

// An Iterator is a generic way of iterating over items.
// Useful for integrating with API that uses pagination.
type Iterator[T any] struct {
	// A function for retrieving the next set of items.
	// If unset the iterator will use the items in the ItemsBuffer directly.
	NextPage NextPageFunc[T]
	// A buffer for storing items. The length controls the size of the buffer.
	// Initial items are used if NextPage is unset.
	ItemsBuffer []T

	err error
}

// Err returns the first error, if any.
func (i *Iterator[T]) Err() error {
	if errors.Is(i.err, errEOI) {
		return nil
	}
	return i.err
}

// Next advances the iterator to the next item.
func (i *Iterator[T]) Next() bool {
	if i.err != nil {
		return false
	}
	if i.NextPage == nil {
		i.NextPage = func([]T) (int, error) { return 0, nil }
		return len(i.ItemsBuffer) > 0
	}
	if len(i.ItemsBuffer) > 1 {
		i.ItemsBuffer = i.ItemsBuffer[1:]
		return true
	}
	if cap(i.ItemsBuffer) == 0 {
		i.ItemsBuffer = make([]T, 1)
	}

	i.ItemsBuffer = i.ItemsBuffer[0:cap(i.ItemsBuffer)]

	count, err := i.NextPage(i.ItemsBuffer)
	if err != nil {
		i.err = fmt.Errorf("next func: %w", err)
		return false
	}
	i.ItemsBuffer = i.ItemsBuffer[:count]

	if len(i.ItemsBuffer) == 0 {
		i.err = errEOI
		return false
	}

	return true
}

// Item is used for retrieving the current item. Callers MUST call i.Next() before calling i.Item().
func (i *Iterator[T]) Item(t *T) {
	*t = i.ItemsBuffer[0]
}

type internalError string

func (e internalError) Error() string { return string(e) }
