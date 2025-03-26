package main

import (
	"context"
	"log/slog"
)

var ctx = context.Background()

func main() {
	serve()
	slog.Info("Started server")

}
