package KCommandParser

import "time"

type Response struct {
	CommandID       string
	Level           string
	Time            time.Time
	ResponseMessage string
}
