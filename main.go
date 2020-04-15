package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"micro-inventory-srv/handler"
	"micro-inventory-srv/subscriber"

	inventory "micro-inventory-srv/proto/inventory"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("mu.micro.book.service.inventory"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	inventory.RegisterInventoryHandler(service.Server(), new(handler.Inventory))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("mu.micro.book.service.inventory", service.Server(), new(subscriber.Inventory))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
