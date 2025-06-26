package main

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Message struct {
	Data string `json:"data"`
}

type PubSub struct {
	subs []chan Message
	mu   sync.Mutex
}

func (ps *PubSub) Subscribe() chan Message {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ch := make(chan Message, 1)
	ps.subs = append(ps.subs, ch)
	return ch
}

func (ps *PubSub) Unsubscribe(ch chan Message) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for i, sub := range ps.subs {
		if sub == ch {
			ps.subs = append(ps.subs[:i], ps.subs[i+1:]...)
			close(ch)
			break
		}
	}
}

func (ps *PubSub) Publish(msg *Message) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	for _, sub := range ps.subs {
		sub <- *msg
	}
}

func main() {
	app := fiber.New()

	pubSub := &PubSub{}

	app.Post("/publisher", func(c *fiber.Ctx) error {
		message := new(Message)
		if err := c.BodyParser(message); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		pubSub.Publish(message)
		return c.JSON(&fiber.Map{
			"message": "Add to subscriber",
		})
	})

	sub := pubSub.Subscribe()
	go func() {
		for msg := range sub {
			fmt.Println("Receive messaage: ", msg)
		}
	}()

	app.Listen(":8080")
}
