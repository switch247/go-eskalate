package controllers

import (
	"bufio"
	"fmt"
	"main/models"
	"main/services"
	"os"
	"strconv"
	"strings"
	"time"
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
	bookID := getIntInput("Enter book ID: ")
	bookTitle := getStringInput("Enter book title: ")
	bookAuthor := getStringInput("Enter book author: ")
	book := &models.Book{
		ID:     bookID,
		Title:  bookTitle,
		Author: bookAuthor,
		Status: "Available",
	}
	err := lc.LibraryService.AddBook(book)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book added successfully", *book)
	}
}

func (lc *LibraryControllerType) AddMember() {
	memberName := getStringInput("Enter member name: ")
	member := &models.Member{
		ID:   len(lc.LibraryService.ListAllMembers()) + 1,
		Name: memberName,
	}
	err := lc.LibraryService.AddMember(member)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Member added successfully", *member)
	}
}

func (lc *LibraryControllerType) RemoveBook() {
	bookID := getIntInput("Enter book ID: ")
	err := lc.LibraryService.RemoveBook(bookID)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book removed successfully")
	}
}

func (lc *LibraryControllerType) BorrowBook() {
	memberID := getIntInput("Enter member ID: ")
	bookID := getIntInput("Enter book ID: ")
	err := lc.LibraryService.BorrowBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book borrowed successfully")
	}
}

func (lc *LibraryControllerType) ReturnBook() {
	memberID := getIntInput("Enter member ID: ")
	bookID := getIntInput("Enter book ID: ")
	err := lc.LibraryService.ReturnBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book returned successfully")
	}
}

func (lc *LibraryControllerType) ListAvailableBooks() {
	books := lc.LibraryService.ListAvailableBooks()
	for _, book := range books {
		fmt.Println(book)
	}
}

func (lc *LibraryControllerType) ListBorrowedBooks() {

	memberID := getIntInput("Enter member ID: ")
	books, err := lc.LibraryService.ListBorrowedBooks(memberID)
	if err == nil {
		for _, book := range books {
			fmt.Println(book)
		}
	} else {
		fmt.Println("Error:", err)
	}
}

func (lc *LibraryControllerType) ListAllMembers() {
	members := lc.LibraryService.ListAllMembers()
	for _, member := range members {
		fmt.Println(member)
	}
}

func (lc *LibraryControllerType) ListAllBooks() {
	books := lc.LibraryService.ListAllBooks()
	for _, book := range books {
		fmt.Println(book)
	}
}

func (lc *LibraryControllerType) ConsoleInteraction() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("================================================================================")
		fmt.Println("Welcome to the Library Management System!")
		fmt.Println("Please select an option:")
		fmt.Println("1. Add a book")
		fmt.Println("2. Remove a book")
		fmt.Println("3. Borrow a book")
		fmt.Println("4. Return a book")
		fmt.Println("5. List available books")
		fmt.Println("6. List borrowed books")
		fmt.Println("7. add a member")
		fmt.Println("8. List all member")
		fmt.Println("9. list all books")
		fmt.Println("0. Exit")
		choice, _ := reader.ReadString('\n')
		choice = string(choice[0])
		switch choice {
		case "1":
			lc.AddBook()
		case "2":
			lc.RemoveBook()
		case "3":
			lc.BorrowBook()
		case "4":
			lc.ReturnBook()
		case "5":
			lc.ListAvailableBooks()
			break
		case "6":
			lc.ListBorrowedBooks()
		case "7":
			lc.AddMember()
		case "8":
			lc.ListAllMembers()
		case "9":
			lc.ListAllBooks()
		case "0":
			fmt.Println("Exiting...")
			reader.ReadString('\n')
			return
		default:
			fmt.Println("Invalid choice")
		}
		fmt.Println("================================================================================")
		reader.ReadString('\n')
		time.Sleep(2000)
	}
}

func getStringInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			return input
		}
		fmt.Println("Invalid input. Please enter a non-empty string.")
	}
}
func getIntInput(prompt string) int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		num, err := strconv.Atoi(input)
		if err == nil {
			return num
		}
		fmt.Println("Invalid input. Please enter an integer.")
	}
}

// int (){
// var test
// }

// list(
// 	LibraryService = int()
// 	def(self) sort(){
// 		self.LibraryService.test
// 	}
// ) //class
// x = list([1,3,4,5,7])
// x.sort()

// LibraryService{
// 	func test(){
// 	}
// }

// controller{
// 	LibraryService = LibraryService()
// 	func (self *controller) chala(){
// 		self.LibraryService.test()
// 	}
// 	LibraryService.test()
// }
