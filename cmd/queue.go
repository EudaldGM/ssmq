package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

type Queue struct {
	name    string
	queue   chan any
	clients []string
}

func newQueue(name string) Queue {
	q := Queue{
		name:  name,
		queue: make(chan any),
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
			slog.Info("Received Message" + msg.(string))
			//write to disk
			q.Send(msg)
		}
	}
}

func (q *Queue) Send(msg any) {
	//for _, client := range q.clients {
	//}
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
}
