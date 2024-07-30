package services

import (
	"errors"
	"main/models"
)

type LibraryService interface {
	AddBook(book *models.Book) error
	RemoveBook(bookID int) error
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) ([]models.Book, error)
}

type libraryService struct {
	books   map[int]*models.Book
	members map[int]*models.Member
}

// AddBook implements LibraryService.
func (l *libraryService) AddBook(book *models.Book) error {
	val := *book
	_, ok := l.books[val.ID]
	if ok {
		return errors.New("Book already exists")
	} else {
		l.books[val.ID] = book
		return nil
	}

}

// BorrowBook implements LibraryService.
func (l *libraryService) BorrowBook(bookID int, memberID int) error {

	book, ok := l.books[bookID]
	if !ok {
		return errors.New("Book not found")
	}

	if book.Status != "Available" {
		return errors.New("Book already borrowed")
	}

	member, ok := l.members[memberID]
	if !ok {
		return errors.New("Member not found")
	}

	newlist := append(member.BorrowedBooks, *(book))
	// test
	l.members[memberID] = &models.Member{
		ID:            member.ID,
		Name:          member.Name,
		BorrowedBooks: newlist,
	}

	l.books[bookID] = &models.Book{
		ID:     book.ID,
		Title:  book.Title,
		Author: book.Author,
		Status: "Borrowed",
	}
	return nil
}

// ListAvailableBooks implements LibraryService.
func (l *libraryService) ListAvailableBooks() []models.Book {
	var availableBooks []models.Book
	for _, book := range l.books {
		if book.Status == "Available" {
			availableBooks = append(availableBooks, *book)
		}
	}
	return availableBooks
}

// ListBorrowedBooks implements LibraryService.
func (l *libraryService) ListBorrowedBooks(memberID int) ([]models.Book, error) {
	member, ok := l.members[memberID]
	if !ok {
		return nil, errors.New("Member not found")
	}
	return member.BorrowedBooks, nil
}

// RemoveBook implements LibraryService.
func (l *libraryService) RemoveBook(bookID int) error {
	_, ok := l.books[bookID]
	if !ok {
		return errors.New("Book not found")
	}
	delete(l.books, bookID)
	return nil
}

// ReturnBook implements LibraryService.
func (l *libraryService) ReturnBook(bookID int, memberID int) error {
	book, ok := l.books[bookID]
	if !ok {
		return errors.New("Book not found")
	}

	if book.Status != "Borrowed" {
		return errors.New("Book not borrowed")
	}

	member, ok := l.members[memberID]
	if !ok {
		return errors.New("Member not found")
	}

	for i, borrowedBook := range member.BorrowedBooks {
		if borrowedBook.ID == bookID {
			deleted_list := append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
			Status := "Available"
			l.books[bookID] = &models.Book{
				ID:     book.ID,
				Title:  book.Title,
				Author: borrowedBook.Author,
				Status: Status,
			}
			l.members[memberID] = &models.Member{
				ID:            member.ID,
				Name:          member.Name,
				BorrowedBooks: deleted_list,
			}
			return nil
		}
	}

	return errors.New("Book not found in member's borrowed books")
}

func NewLibraryService() LibraryService {
	// Dummy data
	_books := map[int]*models.Book{
		1: {ID: 1, Title: "Book 1", Author: "Author 1", Status: "Available"},
		2: {ID: 2, Title: "Book 2", Author: "Author 2", Status: "Available"},
		3: {ID: 3, Title: "Book 3", Author: "Author 3", Status: "Borrowed"},
		4: {ID: 4, Title: "Book 4", Author: "Author 4", Status: "Borrowed"},
		5: {ID: 5, Title: "Book 5", Author: "Author 5", Status: "Available"},
	}

	_members := map[int]*models.Member{
		1: {ID: 1, Name: "Member 1", BorrowedBooks: []models.Book{*(_books[3])}},
		2: {ID: 2, Name: "Member 2", BorrowedBooks: []models.Book{*(_books[4])}},
		3: {ID: 3, Name: "Member 3"},
	}

	return &libraryService{
		books:   _books,
		members: _members,
	}
}
