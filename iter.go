package enrichednow

import (
	"context"
	"encoding/json"
)

// Iter walks a paged list one item at a time, asking the API for the next page when the current one is used up.
//
//	it := client.Runs.Records(ctx, run.ID, nil)
//	for it.Next() {
//		record := it.Current()
//		...
//	}
//	if err := it.Err(); err != nil { ... }
//
// With Go 1.23 or later, range over [Iter.All] instead:
//
//	for record, err := range client.Runs.Records(ctx, run.ID, nil).All() {
//		if err != nil { ... }
//		...
//	}
//
// An Iter is not safe for use by several goroutines at once.
type Iter[T any] struct {
	ctx   context.Context
	fetch func(ctx context.Context, cursor string) (items []T, next string, err error)

	items   []T
	index   int
	current T
	cursor  string
	started bool
	done    bool
	err     error
}

func newIter[T any](ctx context.Context, cursor string, fetch func(ctx context.Context, cursor string) ([]T, string, error)) *Iter[T] {
	return &Iter[T]{ctx: ctx, cursor: cursor, fetch: fetch}
}

// Next advances to the next item and reports whether there is one. It returns false at the end of the list and
// when a page could not be fetched; [Iter.Err] tells the two apart.
func (it *Iter[T]) Next() bool {
	for it.index >= len(it.items) {
		if it.done || it.err != nil {
			return false
		}
		if it.started && it.cursor == "" {
			it.done = true
			return false
		}
		items, next, err := it.fetch(it.ctx, it.cursor)
		it.started = true
		if err != nil {
			it.err = err
			return false
		}
		it.items, it.index, it.cursor = items, 0, next
	}
	it.current = it.items[it.index]
	it.index++
	return true
}

// Current returns the item [Iter.Next] advanced to.
func (it *Iter[T]) Current() T { return it.current }

// Err returns the error that stopped the iteration, if one did.
func (it *Iter[T]) Err() error { return it.err }

// All returns the rest of the list as a Go 1.23 range-over-func iterator (an iter.Seq2). When a page cannot be
// fetched, the last pair carries the error and the zero value.
func (it *Iter[T]) All() func(yield func(T, error) bool) {
	return func(yield func(T, error) bool) {
		for it.Next() {
			if !yield(it.Current(), nil) {
				return
			}
		}
		if it.err != nil {
			var zero T
			yield(zero, it.err)
		}
	}
}

// Collect reads the rest of the list into a slice.
func (it *Iter[T]) Collect() ([]T, error) {
	var out []T
	for it.Next() {
		out = append(out, it.Current())
	}
	return out, it.Err()
}

// Decode decodes one record into T: a struct of yours, or map[string]any.
//
//	type Job struct {
//		Title    string `json:"title"`
//		ApplyURL string `json:"apply_url"`
//	}
//	job, err := enrichednow.Decode[Job](record)
func Decode[T any](record json.RawMessage) (T, error) {
	var out T
	err := json.Unmarshal(record, &out)
	return out, err
}

// DecodeAll reads the rest of a records iterator and decodes every record into T.
func DecodeAll[T any](records *Iter[json.RawMessage]) ([]T, error) {
	var out []T
	for records.Next() {
		item, err := Decode[T](records.Current())
		if err != nil {
			return out, err
		}
		out = append(out, item)
	}
	return out, records.Err()
}
