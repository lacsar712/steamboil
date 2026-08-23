package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lacsar712/steamboil/internal/app"
	"github.com/lacsar712/steamboil/internal/clock"
	"github.com/lacsar712/steamboil/internal/config"
	"github.com/lacsar712/steamboil/internal/model"
)

func TestCase(t *testing.T) {
	clk := clock.NewManual(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	cfg := config.Default("FLAME-1")
	a, err := app.BootstrapWithClock(cfg, clk)
	if err != nil {
		t.Fatal(err)
	}
	comb := a.Snapshot().Combustion
	comb.BurnerPhase = model.BurnerStable
	comb.FurnaceTempF = 400
	if err := a.Store().UpdateCombustion(cfg.UnitID, comb); err != nil {
		t.Fatal(err)
	}
	err = a.OnFlameLoss(context.Background(), "maint-op")
	if err == nil {
		t.Fatal("expected flame loss error")
	}
	if !errors.Is(err, model.ErrFlameLoss) {
		t.Fatalf("expected ErrFlameLoss, got %v", err)
	}
}
