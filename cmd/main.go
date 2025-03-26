package main

import (
	"context"
	"log/slog"
)

var ctx = context.Background()

func init() {
	//initialize from db
}

func main() {
	serve()
	slog.Info("Started server")

}
