package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "root", // no password set
		DB:       0,      // use default DB
	})
}

func setMember(memberId string, duration time.Duration) error {
	err := client.Set(memberId, "active", duration).Err()
	if err != nil {
		return err
	}
	return nil
}

func listenExpired() {
	pubsub := client.Subscribe("__keyevent@0__:expired")
	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

	ch := pubsub.Channel()

	for msg := range ch {
		fmt.Println("expired key:", msg.Payload)
		writeMessage(msg.Payload)
	}
}

func writeMessage(memberId string) {
	// Write a message to your message table
	fmt.Printf("Member %s has expired\n", memberId)
}

func main() {
	go listenExpired()

	err := setMember("member1", 1*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = setMember("member2", 2*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(300 * time.Second)
}
