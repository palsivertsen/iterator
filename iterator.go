package iterator

import (
	"errors"
	"fmt"
)

const (
	errEOI = internalError("end of items")
)

type ItemsFunc[T any] func([]T) (int, error)

type Iterator[T any] struct {
	NextFunc    ItemsFunc[T]
	ItemsBuffer []T

	err error
}

func (i *Iterator[T]) Err() error {
	if errors.Is(i.err, errEOI) {
		return nil
	}
	return i.err
}

func (i *Iterator[T]) Next() bool {
	if i.err != nil {
		return false
	}
	if i.NextFunc == nil {
		i.NextFunc = func([]T) (int, error) { return 0, nil }
		return len(i.ItemsBuffer) > 0
	}
	if len(i.ItemsBuffer) > 1 {
		i.ItemsBuffer = i.ItemsBuffer[1:]
		return true
	}
	if cap(i.ItemsBuffer) == 0 {
		i.ItemsBuffer = make([]T, 32)
	}

	i.ItemsBuffer = i.ItemsBuffer[0:cap(i.ItemsBuffer)]

	count, err := i.NextFunc(i.ItemsBuffer)
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

func (i *Iterator[T]) Item(t *T) {
	*t = i.ItemsBuffer[0]
}

type internalError string

func (e internalError) Error() string { return string(e) }
