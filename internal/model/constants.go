package model

import "time"

const (
	DefaultLeaseTTL        = 30 * time.Second
	PurgeWindow            = 5 * time.Minute
	IgnitionDelayWindow    = 15 * time.Second
	DrumSwellSettleWindow  = 45 * time.Second
	CombustionWarmupWindow = 2 * time.Minute
	FeedwaterRampWindow    = 30 * time.Second
	MaxDrumLevelPercent    = 95.0
	MinDrumLevelPercent    = 15.0
	TripDrumLowPercent     = 10.0
	TripDrumHighPercent    = 98.0
	NormalSteamPressurePSI = 1800.0
	MaxSteamPressurePSI    = 2000.0
	MinFurnaceO2Percent    = 2.5
	MaxFurnaceO2Percent    = 6.0
	DefaultJournalCapacity = 512
)
