package recovery

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"time"
)

type Checkpoint struct {
	ID, FarmID string
	Sequence   int
	At         time.Time
	State      string
}

func New(id, farm, state string, sequence int, at time.Time) (Checkpoint, error) {
	if id == "" || farm == "" || sequence < 0 || state == "" {
		return Checkpoint{}, domain.ErrValidation
	}
	return Checkpoint{ID: id, FarmID: farm, State: state, Sequence: sequence, At: at}, nil
}
func Apply(ctx context.Context, current Checkpoint, next Checkpoint) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if next.FarmID != current.FarmID {
		return domain.ErrForbidden
	}
	if next.Sequence != current.Sequence+1 {
		return fmt.Errorf("%w: sequence", domain.ErrConflict)
	}
	if !next.At.After(current.At) {
		return fmt.Errorf("%w: timestamp", domain.ErrConflict)
	}
	return nil
}
