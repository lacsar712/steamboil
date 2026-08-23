package interlock

import (
	"fmt"

	"github.com/lacsar712/steamboil/internal/model"
)

type PermissiveSet struct {
	fuelOK       bool
	ignitionOK   bool
	drumOK       bool
	pressureOK   bool
	combustionOK bool
}

func NewPermissiveSet() *PermissiveSet { return &PermissiveSet{} }

func (p *PermissiveSet) SetFuel(ok bool)       { p.fuelOK = ok }
func (p *PermissiveSet) SetIgnition(ok bool)   { p.ignitionOK = ok }
func (p *PermissiveSet) SetDrum(ok bool)       { p.drumOK = ok }
func (p *PermissiveSet) SetPressure(ok bool)   { p.pressureOK = ok }
func (p *PermissiveSet) SetCombustion(ok bool) { p.combustionOK = ok }

func (p *PermissiveSet) FuelOK() bool       { return p.fuelOK }
func (p *PermissiveSet) IgnitionOK() bool   { return p.ignitionOK }
func (p *PermissiveSet) DrumOK() bool       { return p.drumOK }
func (p *PermissiveSet) PressureOK() bool   { return p.pressureOK }
func (p *PermissiveSet) CombustionOK() bool { return p.combustionOK }

func (p *PermissiveSet) AllFiring() bool {
	return p.fuelOK && p.ignitionOK && p.drumOK && p.pressureOK && p.combustionOK
}

func (p *PermissiveSet) CheckIgnition() error {
	if !p.fuelOK {
		return fmt.Errorf("%w", model.ErrFuelPermissive)
	}
	if !p.ignitionOK {
		return fmt.Errorf("%w", model.ErrIgnitionBlocked)
	}
	return nil
}

func CheckFlameLoss(reading model.CombustionReading) error {
	if reading.BurnerPhase == model.BurnerStable && reading.FurnaceTempF < 600 {
		return fmt.Errorf("%w", model.ErrFlameLoss)
	}
	return nil
}

func (p *PermissiveSet) CheckFiring() error {
	if err := p.CheckIgnition(); err != nil {
		return err
	}
	if !p.drumOK {
		return fmt.Errorf("%w", model.ErrDrumLevelTrip)
	}
	if !p.pressureOK {
		return fmt.Errorf("%w", model.ErrPressureTrip)
	}
	if !p.combustionOK {
		return fmt.Errorf("%w", model.ErrCombustionTrip)
	}
	return nil
}

type CoordinationLock struct {
	holder string
	held   bool
}

func NewCoordinationLock() *CoordinationLock { return &CoordinationLock{} }

func (c *CoordinationLock) Acquire(holder string) error {
	if c.held {
		return fmt.Errorf("%w", model.ErrCoordinationLock)
	}
	c.holder = holder
	c.held = true
	return nil
}

func (c *CoordinationLock) Release(holder string) {
	if c.held && c.holder == holder {
		c.held = false
		c.holder = ""
	}
}

func (c *CoordinationLock) Require(holder string) error {
	if !c.held || c.holder != holder {
		return fmt.Errorf("%w", model.ErrCoordinationLock)
	}
	return nil
}

func (c *CoordinationLock) Held() bool { return c.held }
