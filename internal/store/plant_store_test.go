package store_test

import (
	"testing"

	"github.com/lacsar712/steamboil/internal/model"
	"github.com/lacsar712/steamboil/internal/store"
)

func TestPlantStoreRegisterGet(t *testing.T) {
	s := store.NewPlantStore()
	snap := model.DefaultSnapshot("U1")
	if err := s.Register("U1", snap); err != nil {
		t.Fatal(err)
	}
	got, ok := s.Get("U1")
	if !ok || got.UnitID != "U1" {
		t.Fatal("expected unit")
	}
}

func TestCompareAndSwapStale(t *testing.T) {
	s := store.NewPlantStore()
	snap := model.DefaultSnapshot("U1")
	_ = s.Register("U1", snap)
	current, _ := s.Get("U1")
	err := s.CompareAndSwap("U1", current.Revision+999, func(ps *model.PlantSnapshot) error {
		ps.State = model.StatePurge
		return nil
	})
	if err == nil {
		t.Fatal("expected stale revision error")
	}
}

func TestJournalAppendRecent(t *testing.T) {
	j := store.NewJournal("", 10)
	for i := 0; i < 15; i++ {
		_, _ = j.Append("U1", "ev", "payload")
	}
	if len(j.Recent(5)) != 5 {
		t.Fatal("expected 5 recent entries")
	}
}
