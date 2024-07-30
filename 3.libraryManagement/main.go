package main

import (
	"main/controllers"
	"main/services"
)

var libraryService services.LibraryService

func main() {
	libraryService = services.NewLibraryService()

	controllers.InitLibraryController(libraryService)
	controllers.LibraryController.ConsoleInteraction()
}
