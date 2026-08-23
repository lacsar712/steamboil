package model_test

import (
	"testing"

	"github.com/lacsar712/steamboil/internal/model"
)

func TestCloneSnapshot(t *testing.T) {
	s := model.DefaultSnapshot("U1")
	s.Alarms = []model.AlarmEvent{{Code: "A1", Active: true}}
	c := model.CloneSnapshot(s)
	c.Alarms[0].Code = "changed"
	if s.Alarms[0].Code != "A1" {
		t.Fatal("clone should deep-copy alarms")
	}
}

func TestDefaultSnapshotLimits(t *testing.T) {
	s := model.DefaultSnapshot("U1")
	if !s.DrumWithinLimits() {
		t.Fatal("default drum should be within limits")
	}
	if s.IsFiring() {
		t.Fatal("cold standby should not be firing")
	}
}

func TestPressureWithinLimits(t *testing.T) {
	s := model.DefaultSnapshot("U1")
	s.State = model.StateFiring
	s.Boiler.SteamPressurePSI = model.MaxSteamPressurePSI + 1
	if s.PressureWithinLimits() {
		t.Fatal("expected trip pressure")
	}
}
