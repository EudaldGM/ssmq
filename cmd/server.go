package main

import (
	"log/slog"
	"net/http"
)

func serve() {
	http.HandleFunc("/newqueue/{queueName}", nq)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func nq(_ http.ResponseWriter, r *http.Request) {
	q := newQueue(r.PathValue("queueName"))
	go q.run(ctx)
	slog.Info("New queue: ", q)
}
