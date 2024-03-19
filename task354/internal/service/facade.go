package service

type Facade struct {
	lib *LibraryService
}

func (f *Facade) UsersList() ([]User, error) {
	users, err := f.lib.GetUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (f *Facade) UserAdd(username string) {

	f.lib.CreateUser(username)
}

func (f *Facade) AuthorsTop(limit int) {
	f.lib.GetTopAuthors(limit)
}

func (f *Facade) AuthorsList() ([]Authors, error) {
	books, err := f.lib.GetAuthorsWithBooks()
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (f *Facade) AuthorAdd(name string) {
	f.lib.AddAuthor(name)
}

func (f *Facade) BookRent(userID int, bookID int) {
	f.lib.RentBook(userID, bookID)
}

func (f *Facade) BookReturn(userID int, bookID int) error {
	err := f.lib.ReturnBook(bookID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (f *Facade) BookList() ([]Book, error) {
	books, err := f.lib.GetAllBooks()
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (f *Facade) BookAdd(title string, authorID int) {
	f.lib.AddBook(title, authorID)
}
