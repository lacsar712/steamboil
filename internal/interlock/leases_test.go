package interlock_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lacsar712/steamboil/internal/interlock"
	"github.com/lacsar712/steamboil/internal/model"
)

func TestLeaseAcquireRelease(t *testing.T) {
	r := interlock.NewLeaseRegistry(time.Minute)
	now := time.Now()
	lease, err := r.Acquire("U1", "op1", now)
	if err != nil {
		t.Fatal(err)
	}
	defer lease.Release()
	if err := r.Require("U1", "op1", now); err != nil {
		t.Fatal(err)
	}
	lease.Release()
	if err := r.Require("U1", "op1", now); !errors.Is(err, model.ErrLeaseMissing) {
		t.Fatalf("expected missing lease, got %v", err)
	}
}

func TestWithLeaseDeferRelease(t *testing.T) {
	r := interlock.NewLeaseRegistry(time.Minute)
	now := time.Now()
	called := false
	err := r.WithLease(context.Background(), "U1", "op1", now, func(ctx context.Context) error {
		called = true
		if err := r.Require("U1", "op2", now); err == nil {
			t.Fatal("other holder should not acquire")
		}
		return nil
	})
	if err != nil || !called {
		t.Fatalf("with lease failed: %v called=%v", err, called)
	}
	if err := r.Require("U1", "op1", now); !errors.Is(err, model.ErrLeaseMissing) {
		t.Fatal("lease should be released after WithLease")
	}
}

func TestPermissiveSet(t *testing.T) {
	p := interlock.NewPermissiveSet()
	p.SetFuel(true)
	p.SetIgnition(true)
	if err := p.CheckIgnition(); err != nil {
		t.Fatal(err)
	}
	p.SetDrum(false)
	if err := p.CheckFiring(); !errors.Is(err, model.ErrDrumLevelTrip) {
		t.Fatalf("got %v", err)
	}
}
