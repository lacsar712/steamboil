package combustion

import (
	"context"
	"fmt"
	"math"

	"github.com/lacsar712/steamboil/internal/clock"
	"github.com/lacsar712/steamboil/internal/model"
)

type Coordinator struct {
	clk     clock.ProcessClock
	burner  *BurnerController
	airflow *AirflowBalancer
	fuel    *FuelRegulator
	purge   *clock.PurgeWindow
	ignition *clock.IgnitionDelayWindow
	warmup  *clock.CombustionWarmupWindow
}

func NewCoordinator(clk clock.ProcessClock) *Coordinator {
	return &Coordinator{
		clk:      clk,
		burner:   NewBurnerController(clk),
		airflow:  NewAirflowBalancer(clk),
		fuel:     NewFuelRegulator(clk),
		purge:    clock.NewPurgeWindow(clk),
		ignition: clock.NewIgnitionDelayWindow(clk),
		warmup:   clock.NewCombustionWarmupWindow(clk),
	}
}

func (c *Coordinator) Burner() *BurnerController  { return c.burner }
func (c *Coordinator) Airflow() *AirflowBalancer { return c.airflow }
func (c *Coordinator) Fuel() *FuelRegulator     { return c.fuel }

func (c *Coordinator) StartPurge(ctx context.Context, snap model.PlantSnapshot) (model.CombustionReading, error) {
	select {
	case <-ctx.Done():
		return snap.Combustion, fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	out := snap.Combustion
	out.BurnerPhase = model.BurnerPurge
	out.PurgeStartedAt = c.clk.Now()
	out.FuelFlowTPH = 0
	out.AirflowTPH = c.airflow.PurgeRate()
	return out, nil
}

func (c *Coordinator) CompletePurge(snap model.CombustionReading) error {
	return c.purge.Require(snap.PurgeStartedAt)
}

func (c *Coordinator) Ignite(ctx context.Context, snap model.PlantSnapshot) (model.CombustionReading, error) {
	select {
	case <-ctx.Done():
		return snap.Combustion, fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	if err := c.purge.Require(snap.Combustion.PurgeStartedAt); err != nil {
		return snap.Combustion, err
	}
	out := snap.Combustion
	out.BurnerPhase = model.BurnerIgnition
	out.IgnitionAt = c.clk.Now()
	out.FuelFlowTPH = c.fuel.IgnitionRate(snap.Settings)
	out.AirflowTPH = c.airflow.IgnitionRate(snap.Settings)
	out.FurnaceTempF = 400
	return out, nil
}

func (c *Coordinator) Stabilize(snap model.PlantSnapshot) (model.CombustionReading, error) {
	if err := c.ignition.Require(snap.Combustion.IgnitionAt); err != nil {
		return snap.Combustion, err
	}
	out := snap.Combustion
	out.BurnerPhase = model.BurnerStable
	out.FuelFlowTPH = snap.Settings.FuelFlowTPH * 0.5
	out.AirflowTPH = c.airflow.Compute(snap)
	out.ExcessO2Pct = c.airflow.ExcessO2(out)
	out.FurnaceTempF = c.burner.EstimateFurnaceTemp(out)
	return out, nil
}

func (c *Coordinator) RampToLoad(snap model.PlantSnapshot, loadPct float64) model.CombustionReading {
	out := snap.Combustion
	out.FuelFlowTPH = snap.Settings.FuelFlowTPH * loadPct
	out.AirflowTPH = c.airflow.Compute(snap)
	out.ExcessO2Pct = c.airflow.ExcessO2(out)
	out.FurnaceTempF = c.burner.EstimateFurnaceTemp(out)
	return out
}

func (c *Coordinator) Trip(snap model.CombustionReading) model.CombustionReading {
	out := snap
	out.BurnerPhase = model.BurnerTrip
	out.FuelFlowTPH = 0
	out.FurnaceTempF = math.Max(200, out.FurnaceTempF*0.5)
	return out
}

func (c *Coordinator) WarmupReady(snap model.CombustionReading) bool {
	return c.warmup.Ready(snap.IgnitionAt)
}
