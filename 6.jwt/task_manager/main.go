package main

import (
	"fmt"
	"log"
	"main/config"
	"main/router"
	// "main/utils"
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
