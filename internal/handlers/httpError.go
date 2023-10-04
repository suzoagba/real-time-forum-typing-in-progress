package handlers

import (
	"fmt"
	"real-time-forum/internal/structs"
)

func HttpError(data structs.ForPage, errorType int, message ...string) structs.ForPage {
	fmt.Println("!!! Error !!!")
	data.HttpError.Error = true
	data.HttpError.Type = errorType
	data.HttpError.Text = message[0]
	if len(message) > 1 {
		data.HttpError.Text2 = message[1]
	}
	return data
}
