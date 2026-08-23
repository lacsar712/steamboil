package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/lacsar712/steamboil/internal/boiler"
	"github.com/lacsar712/steamboil/internal/clock"
	"github.com/lacsar712/steamboil/internal/combustion"
	"github.com/lacsar712/steamboil/internal/config"
	"github.com/lacsar712/steamboil/internal/drum"
	"github.com/lacsar712/steamboil/internal/fsm"
	"github.com/lacsar712/steamboil/internal/interlock"
	"github.com/lacsar712/steamboil/internal/model"
	"github.com/lacsar712/steamboil/internal/store"
)

type App struct {
	cfg           config.Config
	clk           clock.ProcessClock
	store         *store.PlantStore
	journal       *store.Journal
	fsm           *fsm.BoilerFSM
	boiler        *boiler.Controller
	combustion    *combustion.Coordinator
	drum          *drum.Coordinator
	interlock     *interlock.Interlock
	permissives   *interlock.PermissiveSet
	coordLock     *interlock.CoordinationLock
	scheduler     *clock.Scheduler
	purgeWindow   *clock.PurgeWindow
	warmupWindow  *clock.CombustionWarmupWindow
	telemetry     *Telemetry
	tickCancels   map[string]context.CancelFunc
	mu            sync.RWMutex
}

func New(cfg config.Config, clk clock.ProcessClock) *App {
	return &App{
		cfg:          cfg,
		clk:          clk,
		store:        store.NewPlantStore(),
		journal:      store.NewJournal(cfg.JournalPath, cfg.JournalCapacity),
		fsm:          fsm.NewBoilerFSM(cfg.UnitID),
		boiler:       boiler.NewController(clk),
		combustion:   combustion.NewCoordinator(clk),
		drum:         drum.NewCoordinator(clk),
		interlock:    interlock.NewInterlock(cfg.LeaseTTL),
		permissives:  interlock.NewPermissiveSet(),
		coordLock:    interlock.NewCoordinationLock(),
		scheduler:    clock.NewScheduler(clk),
		purgeWindow:  clock.NewPurgeWindow(clk),
		warmupWindow: clock.NewCombustionWarmupWindow(clk),
		telemetry:    NewTelemetry(cfg.UnitID),
		tickCancels:  make(map[string]context.CancelFunc),
	}
}

func (a *App) Snapshot() model.PlantSnapshot {
	snap, err := a.store.Require(a.cfg.UnitID)
	if err != nil {
		return model.DefaultSnapshot(a.cfg.UnitID)
	}
	return snap
}

func (a *App) Config() config.Config              { return a.cfg }
func (a *App) Clock() clock.ProcessClock          { return a.clk }
func (a *App) FSM() *fsm.BoilerFSM                { return a.fsm }
func (a *App) UnitID() string                     { return a.cfg.UnitID }
func (a *App) Store() *store.PlantStore           { return a.store }
func (a *App) Interlock() *interlock.Interlock    { return a.interlock }
func (a *App) Telemetry() TelemetrySnapshot       { return a.telemetry.Snapshot() }
func (a *App) Journal() *store.Journal            { return a.journal }

func (a *App) journalEvent(ev, payload string) {
	_, _ = a.journal.Append(a.cfg.UnitID, ev, payload)
}

func (a *App) syncState(state model.PlantState) {
	_ = a.store.UpdateState(a.cfg.UnitID, state)
}

func (a *App) isFiring(state model.PlantState) bool {
	return state == model.StateFiring || state == model.StateLoadFollow || state == model.StateRamp
}

func (a *App) refreshPermissives(snap model.PlantSnapshot) {
	a.permissives.SetDrum(a.drum.Level().WithinLimits(snap.Drum.LevelPercent))
	a.permissives.SetPressure(a.boiler.Pressure().WithinTripLimits(snap.Boiler.SteamPressurePSI, a.isFiring(snap.State)))
	a.permissives.SetCombustion(a.combustion.Burner().FlameStable(snap.Combustion))
	a.permissives.SetFuel(snap.Combustion.FuelFlowTPH > 0 || snap.State == model.StatePurge)
	a.permissives.SetIgnition(snap.Combustion.BurnerPhase == model.BurnerStable || snap.Combustion.BurnerPhase == model.BurnerIgnition)
	a.fsm.SetFuelPermissive(a.permissives.FuelOK())
	a.fsm.SetPurgeComplete(a.purgeWindow.Ready(snap.Combustion.PurgeStartedAt))
}

func (a *App) tickLabel() string {
	return fmt.Sprintf("%s-tick", a.cfg.UnitID)
}
