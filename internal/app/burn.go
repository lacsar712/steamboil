package app

import (
	"context"
	"time"
)

func (a *App) RunWarmupBurnScheduler(ctx context.Context, ignitionAt time.Time) error {
	for !a.warmupWindow.Ready(ignitionAt) {
		if err := ctx.Err(); err != nil {
			return err
		}
		a.advanceClock(100 * time.Millisecond)
	}
	snap, err := a.store.Require(a.cfg.UnitID)
	if err != nil {
		return err
	}
	// Propagate the caller's context so an operator withdrawing the warmup
	// arrangement cancels the burn-plan installation instead of letting the
	// subsequent ignition steps keep loading into the backend schedule.
	return a.scheduler.InstallBurnPlanCtx(ctx, snap.Settings, "warmup-burn")
}

func (a *App) SchedulerItemCount() int {
	return a.scheduler.ItemCount()
}
