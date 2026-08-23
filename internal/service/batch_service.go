package service

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
)

type BatchResult[T any] struct {
	Successful []T
	Failed     map[string]error
}

func ProcessBatch[T any](ctx context.Context, items []T, key func(T) string, run func(context.Context, T) error) BatchResult[T] {
	result := BatchResult[T]{Successful: []T{}, Failed: map[string]error{}}
	for _, item := range items {
		if err := ctx.Err(); err != nil {
			result.Failed[key(item)] = fmt.Errorf("%w: %v", domain.ErrCancelled, err)
			break
		}
		if err := run(ctx, item); err != nil {
			result.Failed[key(item)] = err
		} else {
			result.Successful = append(result.Successful, item)
		}
	}
	return result
}
