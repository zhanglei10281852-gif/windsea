package ingest

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"sync"
)

type Batch struct {
	ID, FarmID, Source string
	Samples            []domain.Sample
}
type Result struct {
	Accepted, Rejected int
	Errors             []error
}

func ValidateBatch(batch Batch) error {
	if batch.ID == "" || batch.FarmID == "" || batch.Source == "" {
		return domain.ErrValidation
	}
	if len(batch.Samples) == 0 {
		return domain.ErrValidation
	}
	for _, sample := range batch.Samples {
		if err := domain.ValidateSample(sample); err != nil {
			return err
		}
	}
	return nil
}
func Process(ctx context.Context, batch Batch, workers int, fn func(context.Context, domain.Sample) error) Result {
	if workers < 1 {
		workers = 1
	}
	result := Result{Errors: []error{}}
	jobs := make(chan domain.Sample)
	var group sync.WaitGroup
	var mu sync.Mutex
	for i := 0; i < workers; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			for sample := range jobs {
				if err := ctx.Err(); err != nil {
					mu.Lock()
					result.Errors = append(result.Errors, err)
					mu.Unlock()
					continue
				}
				if err := fn(ctx, sample); err != nil {
					mu.Lock()
					result.Rejected++
					result.Errors = append(result.Errors, err)
					mu.Unlock()
				} else {
					mu.Lock()
					result.Accepted++
					mu.Unlock()
				}
			}
		}()
	}
	for _, sample := range batch.Samples {
		select {
		case jobs <- sample:
		case <-ctx.Done():
			break
		}
	}
	close(jobs)
	group.Wait()
	if result.Accepted+result.Rejected < len(batch.Samples) {
		result.Errors = append(result.Errors, fmt.Errorf("%w: processing cancelled", domain.ErrCancelled))
	}
	return result
}
