package gitlab

import "time"

type CommandLog struct {
	Time    time.Time
	Message string
	Error   bool
}

type CommandLogs struct {
	Items []CommandLog
}

func NewCommandLogs() *CommandLogs {
	return &CommandLogs{
		Items: make([]CommandLog, 0),
	}
}

func (l *CommandLogs) Add(
	message string,
	isError bool,
) {
	l.Items = append(
		l.Items,
		CommandLog{
			Time:    time.Now(),
			Message: message,
			Error:   isError,
		},
	)
}

