package drum

import (
	"math"

	"github.com/lacsar712/steamboil/internal/clock"
	"github.com/lacsar712/steamboil/internal/model"
)

type LevelController struct {
	clk clock.ProcessClock
}

func NewLevelController(clk clock.ProcessClock) *LevelController {
	return &LevelController{clk: clk}
}

func (l *LevelController) Compute(snap model.PlantSnapshot, firing bool) (float64, model.DrumCondition) {
	level := snap.Drum.LevelPercent
	if !firing {
		return level, model.DrumNormal
	}
	balance := snap.Drum.FeedwaterTPH - snap.Drum.SteamFlowTPH
	level += balance * 0.01
	level = math.Max(model.MinDrumLevelPercent, math.Min(model.MaxDrumLevelPercent, level))
	cond := l.classify(level, snap)
	return level, cond
}

func (l *LevelController) classify(level float64, snap model.PlantSnapshot) model.DrumCondition {
	setpoint := snap.Settings.DrumLevelSetpoint
	if level > setpoint+15 {
		return model.DrumSwell
	}
	if level < setpoint-15 {
		return model.DrumShrink
	}
	if snap.Boiler.SteamPressurePSI > snap.Settings.TargetSteamPSI*0.9 && level > setpoint+5 {
		return model.DrumCarry
	}
	return model.DrumNormal
}

func (l *LevelController) RecommendFeedwater(snap model.PlantSnapshot, firing bool) float64 {
	if !firing {
		return 0
	}
	err := snap.Settings.DrumLevelSetpoint - snap.Drum.LevelPercent
	return snap.Settings.FeedwaterFlowTPH + err*3
}

func (l *LevelController) WithinLimits(level float64) bool {
	return level >= model.MinDrumLevelPercent && level <= model.MaxDrumLevelPercent
}

func (l *LevelController) TripLow(level float64) bool  { return level < model.TripDrumLowPercent }
func (l *LevelController) TripHigh(level float64) bool { return level > model.TripDrumHighPercent }

func (l *LevelController) LevelError(snap model.PlantSnapshot) float64 {
	return snap.Drum.LevelPercent - snap.Settings.DrumLevelSetpoint
}

func (l *LevelController) ThreeElementBias(snap model.PlantSnapshot) float64 {
	steam := snap.Drum.SteamFlowTPH
	feed := snap.Drum.FeedwaterTPH
	levelErr := l.LevelError(snap)
	return feed + (steam-feed)*0.5 + levelErr*2
}
