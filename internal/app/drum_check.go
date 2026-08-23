package app

import (
	"fmt"

	"github.com/lacsar712/steamboil/internal/model"
)

func (a *App) CheckDrumLevel(snap model.PlantSnapshot) error {
	if snap.Drum.LevelPercent < model.MinDrumLevelPercent {
		return fmt.Errorf("%w", model.ErrDrumLevelLow)
	}
	return nil
}
