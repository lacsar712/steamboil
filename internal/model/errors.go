package model

import "errors"

var (
	ErrContextDone      = errors.New("operation cancelled")
	ErrPlantNotFound    = errors.New("plant unit not found")
	ErrLeaseHeld        = errors.New("interlock lease held by another operator")
	ErrLeaseMissing     = errors.New("interlock lease missing or expired")
	ErrGateBlocked      = errors.New("safety gate blocked")
	ErrFuelPermissive   = errors.New("fuel permissive not satisfied")
	ErrIgnitionBlocked  = errors.New("ignition sequence blocked")
	ErrDrumLevelTrip    = errors.New("drum level trip condition")
	ErrPressureTrip     = errors.New("steam pressure trip condition")
	ErrCombustionTrip   = errors.New("combustion trip condition")
	ErrIllegalState     = errors.New("illegal plant state transition")
	ErrSnapshotStale    = errors.New("snapshot revision stale")
	ErrWindowOpen       = errors.New("timing window still open")
	ErrPurgeIncomplete  = errors.New("furnace purge incomplete")
	ErrCoordinationLock = errors.New("coordination lock held")
)
