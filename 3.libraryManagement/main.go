package main

import (
	"fmt"
	"main/controllers"
	"main/services"
)

var libraryService services.LibraryService

func main() {
	fmt.Println("Hello, World!")
	libraryService = services.NewLibraryService()

	controllers.InitLibraryController(libraryService)
	controllers.LibraryController.ConsoleInteraction()
}
