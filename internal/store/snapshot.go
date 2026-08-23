package store

import "github.com/lacsar712/steamboil/internal/model"

type DrumSnapshotView struct {
	UnitID   string
	Drum     model.DrumReading
	Alarms   []model.AlarmEvent
	Revision uint64
}

func CloneDrumSnapshot(s model.PlantSnapshot) DrumSnapshotView {
	out := DrumSnapshotView{
		UnitID:   s.UnitID,
		Drum:     s.Drum,
		Revision: s.Revision,
	}
	out.Alarms = make([]model.AlarmEvent, len(s.Alarms))
	copy(out.Alarms, s.Alarms)
	return out
}
