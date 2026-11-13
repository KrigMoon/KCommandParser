package KCommandParser

import "time"

type Response struct {
	CommandID       string    `json:"CommandID"`
	Level           string    `json:"Level"`
	Time            time.Time `json:"Time"`
	ResponseMessage string    `json:"ResponseMessage"`
}

func NewResponse(commandID string, level string, responseMessage string) *Response {
	return &Response{
		CommandID:       commandID,
		Time:            time.Now(),
		Level:           level,
		ResponseMessage: responseMessage,
	}
}
