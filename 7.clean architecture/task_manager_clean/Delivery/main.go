package main

import (
	"fmt"
	"log"
	router "main/Delivery/routers"
	"main/config"
)

func init() {
	fmt.Println("Task Manager API")
	_, err := config.MongoInit()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := router.NewRouter()
	if r != nil {
		r.Run()
	} else {
		log.Fatal("Failed to start server")
	}

}
