package app

import (
	"context"
	"fmt"

	"github.com/lacsar712/steamboil/internal/model"
)

const maxBlowdownOpeningPct = 100.0

func (a *App) OpenBlowdown(ctx context.Context, holder string, openingPct float64) error {
	_ = holder
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	if openingPct >= maxBlowdownOpeningPct {
		return fmt.Errorf("blowdown: %w", model.ErrBlowdownLimit)
	}
	return nil
}

func (a *App) BlowdownAfterShutdown(ctx context.Context, openingPct float64) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	snap := a.Snapshot()
	if snap.State != model.StateTrip && snap.State != model.StateColdStandby {
		return fmt.Errorf("plant not shut down")
	}
	if openingPct >= maxBlowdownOpeningPct {
		return fmt.Errorf("unknown fault")
	}
	return nil
}
