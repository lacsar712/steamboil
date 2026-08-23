package combustion_test

import (
	"context"
	"testing"
	"time"

	"github.com/lacsar712/steamboil/internal/clock"
	"github.com/lacsar712/steamboil/internal/combustion"
	"github.com/lacsar712/steamboil/internal/model"
)

func TestStartPurge(t *testing.T) {
	clk := clock.NewManual(time.Now())
	c := combustion.NewCoordinator(clk)
	snap := model.DefaultSnapshot("U1")
	out, err := c.StartPurge(context.Background(), snap)
	if err != nil {
		t.Fatal(err)
	}
	if out.BurnerPhase != model.BurnerPurge {
		t.Fatalf("got %s", out.BurnerPhase)
	}
}

func TestCompletePurgeRequiresWindow(t *testing.T) {
	clk := clock.NewManual(time.Now())
	c := combustion.NewCoordinator(clk)
	snap := model.DefaultSnapshot("U1")
	comb, _ := c.StartPurge(context.Background(), snap)
	if err := c.CompletePurge(comb); err == nil {
		t.Fatal("expected purge incomplete")
	}
	clk.Advance(model.PurgeWindow)
	if err := c.CompletePurge(comb); err != nil {
		t.Fatal(err)
	}
}

func TestAirflowWithinLimits(t *testing.T) {
	clk := clock.NewManual(time.Now())
	a := combustion.NewAirflowBalancer(clk)
	reading := model.CombustionReading{FuelFlowTPH: 10, AirflowTPH: 200, ExcessO2Pct: 3.5}
	if !a.WithinLimits(reading) {
		t.Fatal("expected within limits")
	}
}
