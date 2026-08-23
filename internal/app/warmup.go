package app

import (
	"context"
	"fmt"

	"github.com/lacsar712/steamboil/internal/model"
)

func (a *App) WarmupStatus() (ready bool, detail string) {
	snap := a.Snapshot()
	if snap.Combustion.PurgeStartedAt.IsZero() {
		return false, "purge not started"
	}
	if !a.purgeWindow.Ready(snap.Combustion.PurgeStartedAt) {
		return false, "purge window open"
	}
	if !snap.Combustion.IgnitionAt.IsZero() && !a.warmupWindow.Ready(snap.Combustion.IgnitionAt) {
		return false, "combustion warmup window open"
	}
	if !snap.Drum.LastSwellAt.IsZero() {
		if err := a.drum.RequireSettled(snap.Drum); err != nil {
			return false, "drum swell settling"
		}
	}
	return true, "ready"
}

func (a *App) WaitWarmup(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w", model.ErrContextDone)
		default:
		}
		ready, _ := a.WarmupStatus()
		if ready {
			return nil
		}
	}
}

func (a *App) PurgeRemaining() string {
	snap := a.Snapshot()
	if snap.Combustion.PurgeStartedAt.IsZero() {
		return "not started"
	}
	if a.purgeWindow.Ready(snap.Combustion.PurgeStartedAt) {
		return "complete"
	}
	return "in progress"
}

func (a *App) CombustionWarmupRemaining() string {
	snap := a.Snapshot()
	if snap.Combustion.IgnitionAt.IsZero() {
		return "not ignited"
	}
	if a.warmupWindow.Ready(snap.Combustion.IgnitionAt) {
		return "complete"
	}
	return "in progress"
}
