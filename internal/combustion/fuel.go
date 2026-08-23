package combustion

import (
	"math"

	"github.com/lacsar712/steamboil/internal/clock"
	"github.com/lacsar712/steamboil/internal/model"
)

type FuelRegulator struct {
	clk clock.ProcessClock
}

func NewFuelRegulator(clk clock.ProcessClock) *FuelRegulator {
	return &FuelRegulator{clk: clk}
}

func (f *FuelRegulator) IgnitionRate(settings model.PlantSettings) float64 {
	return settings.FuelFlowTPH * 0.08
}

func (f *FuelRegulator) ComputeForLoad(settings model.PlantSettings, loadPct float64) float64 {
	loadPct = math.Max(0, math.Min(1, loadPct))
	return settings.FuelFlowTPH * loadPct
}

func (f *FuelRegulator) Ramp(current, target, maxStep float64) float64 {
	delta := target - current
	if math.Abs(delta) <= maxStep {
		return target
	}
	if delta > 0 {
		return current + maxStep
	}
	return current - maxStep
}

func (f *FuelRegulator) BtuPerHour(flowTPH float64) float64 {
	return flowTPH * 19_500_000
}

func (f *FuelRegulator) HeatInputMW(flowTPH float64) float64 {
	return flowTPH * 11.6
}

func (f *FuelRegulator) ValidatePermissive(settings model.PlantSettings, drumOK, purgeOK bool) error {
	if !purgeOK {
		return model.ErrPurgeIncomplete
	}
	if !drumOK {
		return model.ErrDrumLevelTrip
	}
	if settings.FuelFlowTPH <= 0 {
		return model.ErrFuelPermissive
	}
	return nil
}

func (f *FuelRegulator) MinFlow(settings model.PlantSettings) float64 {
	return settings.FuelFlowTPH * 0.2
}

func (f *FuelRegulator) MaxFlow(settings model.PlantSettings) float64 {
	return settings.FuelFlowTPH * 1.1
}
