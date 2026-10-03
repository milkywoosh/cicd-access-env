package main

import (
	"access_ctl/internal/config"
	"fmt"
	"log"
)

func main() {
	fmt.Println("test main.go")

	cfg, err := config.LoadConfig("./")
	if err != nil {
		log.Fatalf("err config: %s", err.Error())
	}

	log.Println("RabbitMQUser: ", cfg.RabbitMQUser)
	log.Println("RabbitMQDial: ", cfg.RabbitMQDial)
	log.Println("RabbitMQPassword: ", cfg.RabbitMQPassword)
	log.Println("RabbitMQHost: ", cfg.RabbitMQHost)
}
