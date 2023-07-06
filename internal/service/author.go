package service

import (
	"context"

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

func (a *AuthorService) CreateAuthor(ctx context.Context, req *pb.CreateAuthorRequest) (*pb.AuthorResponse, error) {
	author, err := a.AuthorDB.Create(req.Name)
	if err != nil {
		return nil, err
	}

	authorResponse := pb.AuthorResponse{
		Id:   author.ID,
		Name: author.Name,
	}

	return &authorResponse, nil
}
