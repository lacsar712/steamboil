package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/lacsar712/steamboil/internal/app"
	"github.com/lacsar712/steamboil/internal/clock"
	"github.com/lacsar712/steamboil/internal/config"
	"github.com/lacsar712/steamboil/internal/model"
)

func testApp(t *testing.T) (*app.App, *clock.ManualClock) {
	t.Helper()
	clk := clock.NewManual(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	cfg := config.Default("TEST-1")
	cfg.TickInterval = time.Millisecond
	a, err := app.BootstrapWithClock(cfg, clk)
	if err != nil {
		t.Fatal(err)
	}
	return a, clk
}

func TestBootstrapSnapshot(t *testing.T) {
	a, _ := testApp(t)
	snap := a.Snapshot()
	if snap.UnitID != "TEST-1" {
		t.Fatalf("got %s", snap.UnitID)
	}
	if snap.State != model.StateColdStandby {
		t.Fatalf("got %s", snap.State)
	}
}

func TestStartPurge(t *testing.T) {
	a, _ := testApp(t)
	if err := a.StartPurge(context.Background(), "op"); err != nil {
		t.Fatal(err)
	}
	if a.Snapshot().State != model.StatePurge {
		t.Fatalf("got %s", a.Snapshot().State)
	}
}

func TestTripAndReset(t *testing.T) {
	a, _ := testApp(t)
	if err := a.Trip(context.Background(), "test"); err != nil {
		t.Fatal(err)
	}
	if a.Snapshot().State != model.StateTrip {
		t.Fatal("expected trip")
	}
	if err := a.ResetTrip(context.Background(), "op"); err != nil {
		t.Fatal(err)
	}
	if a.Snapshot().State != model.StateColdStandby {
		t.Fatal("expected cold standby after reset")
	}
}

func TestTickLoop(t *testing.T) {
	a, _ := testApp(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a.StartTickLoop(ctx)
	time.Sleep(20 * time.Millisecond)
	a.StopTickLoop()
	if a.Telemetry().TickCount == 0 {
		t.Fatal("expected ticks")
	}
}
