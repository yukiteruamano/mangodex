package mangodex

import (
	"cmp"
	"context"
	"iter"
)

// Paginator is a public helper for iterating paginated ListResponse[T].
// It is kept public intentionally to avoid duplicating limit/offset/total logic
// across Manga, Chapter, Author, Cover, Group, etc. Thread-safe per instance when used from one goroutine.
type Paginator[T any] struct {
	limit  int
	offset int
	total  int
	fetch  func(ctx context.Context, limit, offset int) (*ListResponse[T], error)
}

// NewPaginator creates a paginator with a fetch function.
func NewPaginator[T any](fetch func(ctx context.Context, limit, offset int) (*ListResponse[T], error)) *Paginator[T] {
	return &Paginator[T]{fetch: fetch, limit: 10}
}

// Next fetches the next page.
func (p *Paginator[T]) Next(ctx context.Context) ([]T, error) {
	if p.fetch == nil {
		return nil, nil
	}
	limit := cmp.Or(p.limit, 10)
	resp, err := p.fetch(ctx, limit, p.offset)
	if err != nil {
		return nil, err
	}
	p.total = resp.Total
	p.offset += len(resp.Data)
	return resp.Data, nil
}

// All returns an iterator over pages (Go 1.23 iter.Seq2).
func (p *Paginator[T]) All(ctx context.Context) iter.Seq2[[]T, error] {
	return func(yield func([]T, error) bool) {
		for p.HasMore() {
			data, err := p.Next(ctx)
			if !yield(data, err) || err != nil {
				return
			}
		}
	}
}

// HasMore reports if more pages are available.
func (p *Paginator[T]) HasMore() bool {
	return p.total == 0 || p.offset < p.total
}

// Reset resets pagination to start.
func (p *Paginator[T]) Reset() {
	p.offset = 0
	p.total = 0
}
