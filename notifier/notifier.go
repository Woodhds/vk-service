package notifier

import (
	"encoding/json"
	"fmt"
	"log"
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

type NotifyService struct {
	messageChan    chan interface{}
	newClients     chan chan interface{}
	closingClients chan chan interface{}
	clients        map[chan interface{}]bool
}

func (n *NotifyService) Listen() {
	for {
		select {
		case s := <-n.newClients:

			// A new client has connected.
			// Register their message channel
			n.clients[s] = true
			log.Printf("Client added. %d registered clients", len(n.clients))
		case s := <-n.closingClients:

			// A client has dettached and we want to
			// stop sending them messages.
			delete(n.clients, s)
			log.Printf("Removed client. %d registered clients", len(n.clients))
		case event := <-n.messageChan:

			// We got a new event from the outside!
			// Send event to all connected clients
			for clientMessageChan := range n.clients {
				clientMessageChan <- event
			}
		}
	}
}

func (n *NotifyService) Success(message string) {
	n.messageChan <- notificationMessage{MessageType: Success, Message: message}
}

func (n *NotifyService) Danger(message string) {
	n.messageChan <- notificationMessage{MessageType: Danger, Message: message}
}

func (n *NotifyService) Warning(message string) {
	n.messageChan <- notificationMessage{MessageType: Warning, Message: message}
}

func NewNotifyService() *NotifyService {
	service := &NotifyService{
		messageChan:    make(chan interface{}, 1),
		newClients:     make(chan chan interface{}),
		closingClients: make(chan chan interface{}),
		clients:        make(map[chan interface{}]bool),
	}

	go service.Listen()

	return service
}

func NotificationHandler(n *NotifyService) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		flusher, _ := rw.(http.Flusher)

		rw.Header().Set("content-type", "text/event-stream")
		rw.Header().Set("cache-control", "no-cache")
		rw.Header().Set("connection", "keep-alive")
		flusher.Flush()

		messageChan := make(chan interface{})
		n.newClients <- messageChan

		defer func() {
			n.closingClients <- messageChan
		}()

		notify := r.Context().Done()

		go func() {
			<-notify
			n.closingClients <- messageChan
		}()

		for {
			data := <-messageChan
			b, _ := json.Marshal(data)
			fmt.Fprintf(rw, "data: %s\n\n", string(b))
			flusher.Flush()
		}
	})
}
