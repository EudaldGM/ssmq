package main

import (
	"context"
	"fmt"
)

type Storage interface {
	Read()
	Write(key string, data string)
}

type RedisStorage struct {
	//Client *redis.Client
	ctx context.Context
}

func (r *RedisStorage) Read() {
	fmt.Println("This is all of the messages")
}

func (r *RedisStorage) Write(key string, data string) {
	fmt.Printf("Message %v has been stored in the queue %v", data, key)
}
