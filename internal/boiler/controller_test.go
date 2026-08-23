package boiler_test

import (
	"context"
	"testing"
	"time"

	"github.com/lacsar712/steamboil/internal/boiler"
	"github.com/lacsar712/steamboil/internal/clock"
	"github.com/lacsar712/steamboil/internal/model"
)

func TestPressureCooldown(t *testing.T) {
	clk := clock.NewManual(time.Now())
	pm := boiler.NewPressureModel(clk)
	p, err := pm.Compute(model.DefaultSnapshot("U1"), false)
	if err != nil {
		t.Fatal(err)
	}
	if p != 0 {
		t.Fatalf("expected 0 pressure when not firing, got %f", p)
	}
}

func TestControllerTick(t *testing.T) {
	clk := clock.NewManual(time.Now())
	c := boiler.NewController(clk)
	snap := model.DefaultSnapshot("U1")
	snap.State = model.StateFiring
	snap.Combustion.FuelFlowTPH = snap.Settings.FuelFlowTPH
	out, err := c.Tick(context.Background(), snap, true)
	if err != nil {
		t.Fatal(err)
	}
	if out.SteamPressurePSI <= 0 {
		t.Fatal("expected positive pressure when firing")
	}
}

func TestRampPressure(t *testing.T) {
	got := boiler.RampPressure(100, 200, 30)
	if got != 130 {
		t.Fatalf("got %f", got)
	}
}
