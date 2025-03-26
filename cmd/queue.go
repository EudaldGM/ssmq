package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Queue struct {
	name    string
	queue   chan io.Reader
	clients []string
}

func newQueue(name string) Queue {
	q := Queue{
		name:  name,
		queue: make(chan io.Reader),
	}
	return q
}

func (q *Queue) run(ctx context.Context) {

	go http.HandleFunc(fmt.Sprintf("/%s", q.name), q.Receive)
	go http.HandleFunc(fmt.Sprintf("/%s/subscribe", q.name), q.Subscribe)
	go http.HandleFunc(fmt.Sprintf("/%s/unsubscribe", q.name), q.Unsubscribe)
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-q.queue:
			//write to disk
			q.Send(msg)
		}
	}
}

func (q *Queue) Send(msg io.Reader) {
	for _, client := range q.clients {
		_, err := http.NewRequest("POST", client, msg)
		if err != nil {
			slog.Error(err.Error())
		}
	}
	slog.Info("Send message to clients")
}

func (q *Queue) Subscribe(w http.ResponseWriter, r *http.Request) {
	q.clients = append(q.clients, r.RemoteAddr)
	slog.Info(fmt.Sprintf("Added client %v to queue %v", r.RemoteAddr, q.name))
	slog.Info("Current clients: ", q.clients)

}

func (q *Queue) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	q.clients = q.clients[:len(q.clients)-1]
	slog.Info(fmt.Sprintf("Added client %v to queue %v", r.RemoteAddr, q.name))
	slog.Info("Current clients: ", q.clients)
}

func (q *Queue) Receive(_ http.ResponseWriter, r *http.Request) {
	q.queue <- r.Body
	slog.Info("Received Message")
}
