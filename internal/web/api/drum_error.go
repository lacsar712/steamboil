package api

import (
	"errors"

	"github.com/lacsar712/steamboil/internal/model"
)

func classifyDrumError(err error) (string, bool) {
	if errors.Is(err, model.ErrDrumLevelLow) {
		return "drum_level_low", true
	}
	return "", false
}
