package zeus_v1_ai

import (
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

func (w *WorkflowsActionsRequest) ConvertWfStrIDs() error {
	for wfi, _ := range w.Workflows {
		if w.Workflows[wfi].WorkflowTemplateStrID != "" {
			wid, werr := strconv.Atoi(w.Workflows[wfi].WorkflowTemplateStrID)
			if werr != nil {
				log.Err(werr).Msg("failed to parse int")
				return werr
			}
			w.Workflows[wfi].WorkflowTemplateID = wid
		}
	}
	return nil
}

func (w *WorkflowsActionsRequest) GetTimeSeriesIterInst() (artemis_orchestrations.Window, bool) {
	if w.UnixStartTime == 0 {
		w.UnixStartTime = int(time.Now().Unix())
	}
	isCycleStepped := false
	endTime := 0
	switch w.DurationUnit {
	case "minutes", "minute":
		minutes := time.Minute
		if w.Duration < 0 {
			w.UnixStartTime += int(minutes.Seconds()) * w.Duration
			w.Duration = -1 * w.Duration
		}
		endTime = w.UnixStartTime + int(minutes.Seconds())*w.Duration
	case "hours", "hour", "hrs", "hr":
		hours := time.Hour
		if w.Duration < 0 {
			w.UnixStartTime += int(hours.Seconds()) * w.Duration
			w.Duration = -1 * w.Duration
		}
		endTime = w.UnixStartTime + int(hours.Seconds())*w.Duration
	case "days", "day":
		days := 24 * time.Hour
		if w.Duration < 0 {
			w.UnixStartTime += int(days.Seconds()) * w.Duration
			w.Duration = -1 * w.Duration
		}
		endTime = w.UnixStartTime + int(days.Seconds())*w.Duration
	case "weeks", "week":
		weeks := 7 * 24 * time.Hour
		if w.Duration < 0 {
			w.UnixStartTime += int(weeks.Seconds()) * w.Duration
			w.Duration = -1 * w.Duration
		}
		endTime = w.UnixStartTime + int(weeks.Seconds())*w.Duration
	case "cycles":
		isCycleStepped = true
	}
	window := artemis_orchestrations.Window{
		UnixStartTime: w.UnixStartTime,
		UnixEndTime:   endTime,
	}
	return window, isCycleStepped
}
