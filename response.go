package KCommandParser

import "time"

// desc: Response structure
type Response struct {
	CommandID       string
	Level           string
	Time            time.Time
	ResponseMessage string
}

// desc: Response constructor
func NewResponse(commandID string, level string, responseMessage string) *Response {
	return &Response{
		CommandID:       commandID,
		Time:            time.Now(),
		Level:           level,
		ResponseMessage: responseMessage,
	}
}
