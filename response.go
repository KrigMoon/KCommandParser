package KCommandParser

import "time"

type Response struct {
	CommandID       string
	Level           string
	Time            time.Time
	ResponseMessage string
}

func NewResponse(commandID string, level string, responseMessage string) *Response {
	return &Response{
		CommandID:       commandID,
		Time:            time.Now(),
		Level:           level,
		ResponseMessage: responseMessage,
	}
}
