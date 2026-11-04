package event

// ErrorPatternMatchedEvent is the first line of a Process's output that matched
// one of its Command's Error Patterns. Later matches in the same Process are not
// reported, so a log storm reads as one event.
type ErrorPatternMatchedEvent struct {
	CommandId string
	Line      string
}

func (ErrorPatternMatchedEvent) GetName() string {
	return "domain_event.process.error_pattern_matched"
}
