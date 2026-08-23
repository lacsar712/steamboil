package clock_test

import (
	"testing"
	"time"

	"github.com/lacsar712/steamboil/internal/clock"
	"github.com/lacsar712/steamboil/internal/model"
)

func TestManualClockAdvance(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	c := clock.NewManual(start)
	c.Advance(time.Minute)
	if c.Now().Sub(start) != time.Minute {
		t.Fatalf("expected 1 minute advance")
	}
}

func TestPurgeWindow(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	c := clock.NewManual(start)
	p := clock.NewPurgeWindow(c)
	if p.Ready(start) {
		t.Fatal("purge should not be ready immediately")
	}
	c.Advance(model.PurgeWindow)
	if !p.Ready(start) {
		t.Fatal("purge should be ready after window")
	}
}

func TestFeedwaterRampTracker(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	c := clock.NewManual(start)
	tr := clock.NewFeedwaterRampTracker(c)
	if tr.Satisfied() {
		t.Fatal("ramp should not be satisfied immediately")
	}
	c.Advance(model.FeedwaterRampWindow)
	if !tr.Satisfied() {
		t.Fatal("ramp should be satisfied after window")
	}
}
