package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
)

type client struct {
	ip net.IP
	//port int
}

func newClient(ip net.IP /*port int*/) (client, error) {
	//implement correct json parsing for the client
	//if port < 1 || port > 65535 {
	//	return client{}, fmt.Errorf("invalid port %d", port)
	//}
	sqc := client{ip: ip /*port: port*/}
	return sqc, nil
}

type Queue struct {
	name    string
	queue   chan io.Reader
	clients []client
}

func newQueue(name string) Queue {
	q := Queue{
		name:    name,
		queue:   make(chan io.Reader),
		clients: []client{},
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
		_, err := http.NewRequest("POST", client.ip.String(), msg)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
	}
	slog.Info("Message sent to clients")
}

func (q *Queue) Subscribe(w http.ResponseWriter, r *http.Request) {
	//implement correct json parsing for the client
	ip, _, _ := net.ParseCIDR(r.RemoteAddr)
	c, err := newClient(ip)
	if err != nil {
		slog.Error(err.Error())
	}
	q.clients = append(q.clients, c)
	slog.Info(fmt.Sprintf("Added client %v to queue %v", r.RemoteAddr, q.name))
}

func (q *Queue) Unsubscribe(_ http.ResponseWriter, r *http.Request) {
	//add removal logic
	slog.Info(fmt.Sprintf("Removed client %v from queue %v", r.RemoteAddr, q.name))
}

func (q *Queue) Receive(_ http.ResponseWriter, r *http.Request) {
	q.queue <- r.Body
	slog.Info("Received Message")
}
