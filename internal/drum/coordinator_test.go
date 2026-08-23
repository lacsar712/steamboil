package drum_test

import (
	"context"
	"testing"
	"time"

	"github.com/lacsar712/steamboil/internal/clock"
	"github.com/lacsar712/steamboil/internal/drum"
	"github.com/lacsar712/steamboil/internal/model"
)

func TestDrumTickFiring(t *testing.T) {
	clk := clock.NewManual(time.Now())
	c := drum.NewCoordinator(clk)
	snap := model.DefaultSnapshot("U1")
	snap.State = model.StateFiring
	snap.Boiler.MainSteamFlowTPH = 100
	out, err := c.Tick(context.Background(), snap, true)
	if err != nil {
		t.Fatal(err)
	}
	if out.SteamFlowTPH <= 0 {
		t.Fatal("expected steam flow when firing")
	}
}

func TestTripRequired(t *testing.T) {
	c := drum.NewCoordinator(clock.NewManual(time.Now()))
	if !c.TripRequired(model.DrumReading{LevelPercent: 5}) {
		t.Fatal("expected low drum trip")
	}
}

func TestSwellDetector(t *testing.T) {
	d := drum.NewSwellDetector(clock.NewManual(time.Now()))
	_, _ = d.Observe(50)
	changed, cond := d.Observe(50)
	if changed || cond != model.DrumNormal {
		t.Fatal("steady level should be normal")
	}
	changed, cond = d.Observe(60)
	if !changed || cond != model.DrumSwell {
		t.Fatalf("expected swell, got %v %s", changed, cond)
	}
}
