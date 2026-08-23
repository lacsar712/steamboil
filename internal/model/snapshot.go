package model

import "time"

func CloneSnapshot(s PlantSnapshot) PlantSnapshot {
	out := s
	out.Alarms = append([]AlarmEvent(nil), s.Alarms...)
	return out
}

func DefaultSnapshot(unitID string) PlantSnapshot {
	now := time.Now()
	return PlantSnapshot{
		UnitID: unitID,
		State:  StateColdStandby,
		Settings: PlantSettings{
			Mode:              ModeBaseLoad,
			TargetMW:          150,
			TargetSteamPSI:    NormalSteamPressurePSI,
			DrumLevelSetpoint: 55,
			FeedwaterFlowTPH:  400,
			FuelFlowTPH:       35,
			ExcessO2Setpoint:  3.5,
		},
		Plant: PlantRef{UnitLabel: unitID, PlantCode: "STEAM-PLT"},
		Drum: DrumReading{
			LevelPercent: 50,
			Condition:    DrumNormal,
			FeedwaterTPH: 0,
			SteamFlowTPH: 0,
		},
		Combustion: CombustionReading{
			BurnerPhase: BurnerIdle,
		},
		Boiler: BoilerReading{
			SteamPressurePSI: 0,
			SteamTempF:       70,
		},
		UpdatedAt: now,
	}
}

func (s PlantSnapshot) IsFiring() bool {
	return s.State == StateFiring || s.State == StateLoadFollow || s.State == StateRamp
}

func (s PlantSnapshot) DrumWithinLimits() bool {
	return s.Drum.LevelPercent >= MinDrumLevelPercent && s.Drum.LevelPercent <= MaxDrumLevelPercent
}

func (s PlantSnapshot) PressureWithinLimits() bool {
	if !s.IsFiring() {
		return true
	}
	return s.Boiler.SteamPressurePSI <= MaxSteamPressurePSI
}
