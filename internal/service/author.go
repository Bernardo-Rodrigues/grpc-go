package service

import (
	"context"
	"io"

	"github.com/Bernardo-Rodrigues/grpc-go/internal/database"
	"github.com/Bernardo-Rodrigues/grpc-go/internal/pb"
)

type AuthorService struct {
	pb.UnimplementedAuthorServiceServer
	AuthorDB database.Author
}

func NewAuthorService(authorDB database.Author) *AuthorService {
	return &AuthorService{
		AuthorDB: authorDB,
	}
}

func (a *AuthorService) CreateAuthor(ctx context.Context, req *pb.CreateAuthorRequest) (*pb.Author, error) {
	author, err := a.AuthorDB.Create(req.Name)
	if err != nil {
		return nil, err
	}

	authorResponse := pb.Author{
		Id:   author.ID,
		Name: author.Name,
	}

	return &authorResponse, nil
}

func (a *AuthorService) CreateAuthorStream(stream pb.AuthorService_CreateAuthorStreamServer) error {
	authors := &pb.AuthorList{}

	for {
		authorRequest, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(authors)
		}
		if err != nil {
			return err
		}

		author, err := a.AuthorDB.Create(authorRequest.Name)
		if err != nil {
			return err
		}

		authors.Authors = append(authors.Authors, &pb.Author{
			Id:   author.ID,
			Name: author.Name,
		})
	}
}

func (a *AuthorService) ListAuthors(ctx context.Context, req *pb.Blank) (*pb.AuthorList, error) {
	authors, err := a.AuthorDB.FindAll()
	if err != nil {
		return nil, err
	}

	var authorsList []*pb.Author

	for _, author := range authors {
		authorResponse := &pb.Author{
			Id:   author.ID,
			Name: author.Name,
		}

		authorsList = append(authorsList, authorResponse)
	}

	return &pb.AuthorList{Authors: authorsList}, nil
}

func (a *AuthorService) GetAuthor(ctx context.Context, req *pb.GetAuthorRequest) (*pb.Author, error) {
	author, err := a.AuthorDB.Find(req.Id)
	if err != nil {
		return nil, err
	}

	authorResponse := pb.Author{
		Id:   author.ID,
		Name: author.Name,
	}

	return &authorResponse, nil
}
