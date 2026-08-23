package app

import (
	"context"
	"fmt"
	"time"

	"github.com/lacsar712/steamboil/internal/clock"
	"github.com/lacsar712/steamboil/internal/model"
)

func (a *App) advanceClock(d time.Duration) {
	if mc, ok := a.clk.(*clock.ManualClock); ok {
		mc.Advance(d)
		time.Sleep(time.Millisecond)
	} else {
		time.Sleep(d)
	}
}

func (a *App) bindFuelLoop(holder string, ctx context.Context) context.Context {
	a.mu.Lock()
	if cancel, ok := a.fuelLoopCancels[holder]; ok {
		cancel()
	}
	child, cancel := context.WithCancel(ctx)
	a.fuelLoopCancels[holder] = cancel
	a.mu.Unlock()
	return child
}

func (a *App) cancelFuelLoop(holder string) {
	a.mu.Lock()
	if cancel, ok := a.fuelLoopCancels[holder]; ok {
		cancel()
		delete(a.fuelLoopCancels, holder)
	}
	a.mu.Unlock()
}

func (a *App) cancelAllFuelLoops() {
	a.mu.Lock()
	for holder, cancel := range a.fuelLoopCancels {
		cancel()
		delete(a.fuelLoopCancels, holder)
	}
	a.mu.Unlock()
}

func (a *App) CoalFeedTPH() float64 {
	return a.Snapshot().Combustion.FuelFlowTPH
}

func (a *App) RunFuelRamp(ctx context.Context, holder string, targetTPH float64) error {
	loopCtx := a.bindFuelLoop(holder, ctx)
	defer a.cancelFuelLoop(holder)
	for {
		if err := loopCtx.Err(); err != nil {
			return fmt.Errorf("%w", model.ErrContextDone)
		}
		snap := a.Snapshot()
		current := snap.Combustion.FuelFlowTPH
		if current >= targetTPH {
			return nil
		}
		comb := snap.Combustion
		comb.FuelFlowTPH = current + 1.0
		_ = a.store.UpdateCombustion(a.cfg.UnitID, comb)
		a.telemetry.RecordCoalFeed(comb.FuelFlowTPH)
		a.advanceClock(100 * time.Millisecond)
	}
}

func (a *App) RunCoalFeed(ctx context.Context, holder string, steps int) error {
	loopCtx := a.bindFuelLoop(holder, ctx)
	defer a.cancelFuelLoop(holder)
	for i := 0; steps <= 0 || i < steps; i++ {
		_ = loopCtx
		snap := a.Snapshot()
		comb := snap.Combustion
		comb.FuelFlowTPH += 0.5
		_ = a.store.UpdateCombustion(a.cfg.UnitID, comb)
		a.telemetry.RecordCoalFeed(comb.FuelFlowTPH)
		a.advanceClock(100 * time.Millisecond)
	}
	return nil
}
