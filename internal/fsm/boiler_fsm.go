package fsm

import (
	"context"
	"fmt"
	"sync"

	"github.com/lacsar712/steamboil/internal/model"
)

type BoilerFSM struct {
	mu            sync.RWMutex
	state         model.PlantState
	fuelPermissive bool
	purgeComplete  bool
	hooks          *HookChain
}

func NewBoilerFSM(unitID string) *BoilerFSM {
	_ = unitID
	return &BoilerFSM{state: model.StateColdStandby, hooks: NewHookChain()}
}

func (f *BoilerFSM) Hooks() *HookChain { return f.hooks }

func (f *BoilerFSM) State() model.PlantState {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.state
}

func (f *BoilerFSM) SetFuelPermissive(ok bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fuelPermissive = ok
}

func (f *BoilerFSM) SetPurgeComplete(ok bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.purgeComplete = ok
}

func (f *BoilerFSM) FuelPermissive() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.fuelPermissive
}

func (f *BoilerFSM) Dispatch(ctx context.Context, event PlantEvent) (model.PlantState, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	select {
	case <-ctx.Done():
		return f.state, fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	if event == EvTrip {
		from := f.state
		if f.hooks != nil {
			if err := f.hooks.RunBefore(ctx, from, model.StateTrip, event); err != nil {
				return f.state, err
			}
		}
		f.state = model.StateTrip
		if f.hooks != nil {
			if err := f.hooks.RunAfter(ctx, from, model.StateTrip, event); err != nil {
				return f.state, err
			}
		}
		return f.state, nil
	}
	next, ok := NextState(f.state, event)
	if !ok {
		if f.hooks != nil {
			_ = f.hooks.RunAfter(ctx, f.state, f.state, event)
		}
		return f.state, fmt.Errorf("%s from %s: %w", event, f.state, ErrIllegalTransition)
	}
	if event == EvIgnite && !f.fuelPermissive {
		return f.state, fmt.Errorf("%w", model.ErrFuelPermissive)
	}
	if event == EvPurgeComplete && !f.purgeComplete {
		return f.state, fmt.Errorf("%w", model.ErrPurgeIncomplete)
	}
	from := f.state
	if f.hooks != nil {
		if err := f.hooks.RunBefore(ctx, from, next, event); err != nil {
			return f.state, err
		}
	}
	f.state = next
	if f.hooks != nil {
		if err := f.hooks.RunAfter(ctx, from, next, event); err != nil {
			return f.state, err
		}
	}
	return f.state, nil
}

func (f *BoilerFSM) ForceState(state model.PlantState) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.state = state
}
