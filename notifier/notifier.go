package notifier

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	Success = iota
	Danger  = iota
	Warning = iota
)

type notificationMessage struct {
	MessageType int    `json:"messageType"`
	Message     string `json:"message"`
}

type Notifier interface {
	Listen() http.Handler
	Success(message string)
	Danger(message string)
	Warning(message string)
}

type notifyService struct {
	messageChan chan notificationMessage
}

func (n notifyService) Listen() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		flusher, _ := rw.(http.Flusher)

		rw.Header().Set("content-type", "text/event-stream")
		rw.Header().Set("cache-control", "no-cache")
		rw.Header().Set("connection", "keep-alive")
		flusher.Flush()

		for {
			select {
			case data := <-n.messageChan:
				d, _ := json.Marshal(data)
				fmt.Fprintf(rw, "data: %s\n\n", string(d))
				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}
	})
}

func (n notifyService) Success(message string) {
	n.messageChan <- notificationMessage{MessageType: Success, Message: message}
}

func (n notifyService) Danger(message string) {
	n.messageChan <- notificationMessage{MessageType: Danger, Message: message}
}

func (n notifyService) Warning(message string) {
	n.messageChan <- notificationMessage{MessageType: Warning, Message: message}
}

func NewNotifyService() Notifier {
	return &notifyService{
		messageChan: make(chan notificationMessage),
	}
}
