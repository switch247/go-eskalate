package controllers

import (
	"bufio"
	"fmt"
	"main/models"
	"main/services"
	"os"
	"strconv"
	"strings"
)

type LibraryControllerType struct {
	LibraryService services.LibraryService
}

var LibraryController *LibraryControllerType

func InitLibraryController(libraryService services.LibraryService) {
	LibraryController = &LibraryControllerType{
		LibraryService: libraryService,
	}
}

func (lc *LibraryControllerType) AddBook() {
	lc.LibraryService.AddBook(&models.Book{})
}

func (lc *LibraryControllerType) RemoveBook(bookID int) {
	lc.LibraryService.RemoveBook(bookID)
}

func (lc *LibraryControllerType) BorrowBook(bookID int, memberID int) {
	lc.LibraryService.BorrowBook(bookID, memberID)
}

func (lc *LibraryControllerType) ReturnBook(bookID int, memberID int) {
	lc.LibraryService.ReturnBook(bookID, memberID)
}

func (lc *LibraryControllerType) ListAvailableBooks() {
	books := lc.LibraryService.ListAvailableBooks()
	for _, book := range books {
		fmt.Println(book)
	}
}

func (lc *LibraryControllerType) ListBorrowedBooks() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter member ID: ")
	memberIDString, _ := reader.ReadString('\n')
	memberID, err := strconv.Atoi(strings.TrimSpace(memberIDString))
	if err != nil {
		fmt.Println("Invalid member ID")
		return
	}
	books, err := lc.LibraryService.ListBorrowedBooks(memberID)
	if err == nil {
		for _, book := range books {
			fmt.Println(book)
		}
	} else {
		fmt.Println("Error:", err)
	}
	return
}

func (lc *LibraryControllerType) ConsoleInteraction() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Library Management System!")
	fmt.Println("Please select an option:")
	fmt.Println("1. Add a book")
	fmt.Println("2. Remove a book")
	fmt.Println("3. Borrow a book")
	fmt.Println("4. Return a book")
	fmt.Println("5. List available books")
	fmt.Println("6. List borrowed books")
	fmt.Println("7. Exit")

	for {
		choice, _ := reader.ReadString('\n')
		choice = string(choice[0])
		switch choice {
		case "1":
			lc.AddBook()
		case "2":
			lc.RemoveBook(1)
		case "3":
			lc.BorrowBook(1, 1)
		case "4":
			lc.ReturnBook(1, 1)
		case "5":
			lc.ListAvailableBooks()
			break
		case "6":
			lc.ListBorrowedBooks()
		case "7":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice")
			// lc.ConsoleInteraction()
		}
	}
}

func getIntInput() int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter input: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		num, err := strconv.Atoi(input)
		if err == nil {
			return num
		}
		fmt.Println("Invalid input. Please enter an integer.")
	}
}
