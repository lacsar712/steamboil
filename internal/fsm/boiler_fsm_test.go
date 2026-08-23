package fsm_test

import (
	"context"
	"testing"

	"github.com/lacsar712/steamboil/internal/fsm"
	"github.com/lacsar712/steamboil/internal/model"
)

func TestColdStandbyToPurge(t *testing.T) {
	f := fsm.NewBoilerFSM("U1")
	state, err := f.Dispatch(context.Background(), fsm.EvStartPurge)
	if err != nil {
		t.Fatal(err)
	}
	if state != model.StatePurge {
		t.Fatalf("got %s", state)
	}
}

func TestIgniteRequiresFuelPermissive(t *testing.T) {
	f := fsm.NewBoilerFSM("U1")
	f.SetFuelPermissive(false)
	f.SetPurgeComplete(true)
	_, _ = f.Dispatch(context.Background(), fsm.EvStartPurge)
	_, _ = f.Dispatch(context.Background(), fsm.EvPurgeComplete)
	_, err := f.Dispatch(context.Background(), fsm.EvIgnite)
	if err == nil {
		t.Fatal("expected fuel permissive error")
	}
}

func TestContextCancel(t *testing.T) {
	f := fsm.NewBoilerFSM("U1")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := f.Dispatch(ctx, fsm.EvStartPurge)
	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestIllegalTransition(t *testing.T) {
	f := fsm.NewBoilerFSM("U1")
	_, err := f.Dispatch(context.Background(), fsm.EvReachFiring)
	if err == nil {
		t.Fatal("expected illegal transition")
	}
}
