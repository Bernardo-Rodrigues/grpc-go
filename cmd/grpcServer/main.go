package main

import (
	"database/sql"
	"net"

	"github.com/Bernardo-Rodrigues/grpc-go/internal/database"
	"github.com/Bernardo-Rodrigues/grpc-go/internal/pb"
	"github.com/Bernardo-Rodrigues/grpc-go/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	authorDb := database.NewAuthor(db)
	authorService := service.NewAuthorService(*authorDb)
	bookDb := database.NewBook(db)
	bookService := service.NewBookService(*bookDb)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthorServiceServer(grpcServer, authorService)
	pb.RegisterBookServiceServer(grpcServer, bookService)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
