package services

import (
	"main/models"
)

type LibraryService interface {
	AddBook(book *models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

type libraryService struct {
	books   []models.Book
	members []models.Member
}

// AddBook implements LibraryService.
func (l *libraryService) AddBook(book *models.Book) {
	panic("unimplemented")
}

// BorrowBook implements LibraryService.
func (l *libraryService) BorrowBook(bookID int, memberID int) error {
	panic("unimplemented")
}

// ListAvailableBooks implements LibraryService.
func (l *libraryService) ListAvailableBooks() []models.Book {
	// panic("unimplemented")

	var availableBooks []models.Book
	for _, book := range l.books {
		if book.Status == "Available" {
			availableBooks = append(availableBooks, book)
		}
	}
	return availableBooks
}

// ListBorrowedBooks implements LibraryService.
func (l *libraryService) ListBorrowedBooks(memberID int) []models.Book {
	panic("unimplemented")
}

// RemoveBook implements LibraryService.
func (l *libraryService) RemoveBook(bookID int) {
	panic("unimplemented")
}

// ReturnBook implements LibraryService.
func (l *libraryService) ReturnBook(bookID int, memberID int) error {
	panic("unimplemented")
}

func NewLibraryService() LibraryService {
	//  dummy data
	var _books = []models.Book{
		{ID: 1, Title: "Book 1", Author: "Author 1", Status: "Available"},
		{ID: 2, Title: "Book 2", Author: "Author 2", Status: "Available"},
		{ID: 3, Title: "Book 3", Author: "Author 3", Status: "Borrowed"},
		{ID: 4, Title: "Book 4", Author: "Author 4", Status: "Borrowed"},
		{ID: 5, Title: "Book 5", Author: "Author 5", Status: "Available"},
	}

	var _members = []models.Member{
		{ID: 1, Name: "Member 1", BorrowedBooks: []models.Book{
			_books[2], // Boook3 is borrowed by member 1
		}},
		{ID: 2, Name: "Member 2", BorrowedBooks: []models.Book{
			_books[3],
		}},
		{ID: 3, Name: "Member 3"},
	}

	return &libraryService{
		// books:   make([]models.Book, 0),
		// members: make([]models.Member, 0),
		books:   _books,
		members: _members,
	}
}
