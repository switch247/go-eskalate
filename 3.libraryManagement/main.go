package main

import (
	"fmt"
	"main/services"
)

var libraryService services.LibraryService

func main() {
	libraryService = services.NewLibraryService()
	fmt.Println("Hello, World!")
	fmt.Printf("%+v\n", libraryService.ListAvailableBooks())
}
