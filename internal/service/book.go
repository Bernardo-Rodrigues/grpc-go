package service

import (
	"context"
	"io"

	"github.com/Bernardo-Rodrigues/grpc-go/internal/database"
	"github.com/Bernardo-Rodrigues/grpc-go/internal/pb"
)

type BookService struct {
	pb.UnimplementedBookServiceServer
	BookDB database.Book
}

func NewBookService(bookDB database.Book) *BookService {
	return &BookService{
		BookDB: bookDB,
	}
}

func (a *BookService) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.Book, error) {
	book, err := a.BookDB.Create(req.Name, req.Description, req.AuthorId)
	if err != nil {
		return nil, err
	}

	bookResponse := pb.Book{
		Id:          book.ID,
		Name:        book.Name,
		Description: book.Description,
		AuthorId:    book.AuthorID,
	}

	return &bookResponse, nil
}

func (a *BookService) CreateBookStream(stream pb.BookService_CreateBookStreamServer) error {
	books := &pb.BookList{}

	for {
		bookRequest, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(books)
		}
		if err != nil {
			return err
		}

		book, err := a.BookDB.Create(bookRequest.Name, bookRequest.Description, bookRequest.AuthorId)
		if err != nil {
			return err
		}

		books.Books = append(books.Books, &pb.Book{
			Id:          book.ID,
			Name:        book.Name,
			Description: book.Description,
			AuthorId:    book.AuthorID,
		})
	}
}

func (a *BookService) CreateBookStreamBidirectional(stream pb.BookService_CreateBookStreamBidirectionalServer) error {
	for {
		bookRequest, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		book, err := a.BookDB.Create(bookRequest.Name, bookRequest.Description, bookRequest.AuthorId)
		if err != nil {
			return err
		}

		err = stream.Send(&pb.Book{
			Id:          book.ID,
			Name:        book.Name,
			Description: book.Description,
			AuthorId:    book.AuthorID,
		})

		if err != nil {
			return err
		}
	}
}

func (a *BookService) ListBooks(ctx context.Context, req *pb.Blank) (*pb.BookList, error) {
	books, err := a.BookDB.FindAll()
	if err != nil {
		return nil, err
	}

	var booksList []*pb.Book

	for _, book := range books {
		bookResponse := &pb.Book{
			Id:          book.ID,
			Name:        book.Name,
			Description: book.Description,
			AuthorId:    book.AuthorID,
		}

		booksList = append(booksList, bookResponse)
	}

	return &pb.BookList{Books: booksList}, nil
}

func (a *BookService) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.Book, error) {
	book, err := a.BookDB.Find(req.Id)
	if err != nil {
		return nil, err
	}

	bookResponse := pb.Book{
		Id:   book.ID,
		Name: book.Name,
	}

	return &bookResponse, nil
}
