package apollo_pagerduty

const (
	INFO     = "info"
	WARNING  = "warning"
	ERROR    = "error"
	CRITICAL = "critical"

	TRIGGER     = "trigger"
	ACKNOWLEDGE = "acknowledge"
	RESOLVE     = "resolve"
)

var (
	Sev    Severity
	Action EventAction
)

type Severity string

func (s *Severity) Info() string {
	return INFO
}
func (s *Severity) Error() string {
	return ERROR
}
func (s *Severity) Warning() string {
	return WARNING
}
func (s *Severity) Critical() string {
	return CRITICAL
}

type EventAction string

func (e *EventAction) Trigger() string {
	return TRIGGER
}
func (e *EventAction) Acknowledge() string {
	return ACKNOWLEDGE
}
func (e *EventAction) Resolve() string {
	return RESOLVE
}
