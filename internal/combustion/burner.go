package combustion

import (
	"math"

	"github.com/lacsar712/steamboil/internal/clock"
	"github.com/lacsar712/steamboil/internal/model"
)

type BurnerController struct {
	clk clock.ProcessClock
}

func NewBurnerController(clk clock.ProcessClock) *BurnerController {
	return &BurnerController{clk: clk}
}

func (b *BurnerController) EstimateFurnaceTemp(reading model.CombustionReading) float64 {
	base := 300.0
	fuelHeat := reading.FuelFlowTPH * 50
	airCool := reading.AirflowTPH * 2
	return base + fuelHeat - airCool
}

func (b *BurnerController) FlameStable(reading model.CombustionReading) bool {
	if reading.BurnerPhase != model.BurnerStable && reading.BurnerPhase != model.BurnerIgnition {
		return false
	}
	return reading.FurnaceTempF > 800 && reading.ExcessO2Pct >= model.MinFurnaceO2Percent
}

func (b *BurnerController) TripRequired(reading model.CombustionReading) bool {
	if reading.ExcessO2Pct > model.MaxFurnaceO2Percent*2 {
		return true
	}
	if reading.BurnerPhase == model.BurnerTrip {
		return true
	}
	if reading.FurnaceTempF > 3500 {
		return true
	}
	return false
}

func (b *BurnerController) PhaseLabel(phase model.BurnerPhase) string {
	switch phase {
	case model.BurnerIdle:
		return "Idle"
	case model.BurnerPurge:
		return "Purge"
	case model.BurnerIgnition:
		return "Ignition"
	case model.BurnerStable:
		return "Stable Flame"
	case model.BurnerTrip:
		return "Tripped"
	default:
		return string(phase)
	}
}

func (b *BurnerController) HeatReleaseMW(reading model.CombustionReading) float64 {
	return reading.FuelFlowTPH * 12.5
}

func (b *BurnerController) TurndownRatio(settings model.PlantSettings, currentFuel float64) float64 {
	if settings.FuelFlowTPH <= 0 {
		return 0
	}
	return currentFuel / settings.FuelFlowTPH
}

func (b *BurnerController) MinStableFuel(settings model.PlantSettings) float64 {
	return settings.FuelFlowTPH * 0.25
}

func (b *BurnerController) NormalizeFuel(flow, max float64) float64 {
	return math.Min(math.Max(flow, 0), max)
}
