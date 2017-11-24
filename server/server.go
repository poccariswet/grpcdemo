package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	pb "github.com/soeyusuke/grpcdemo/proto"
)

const (
	port = ":1111"
)

var (
	tablename = "bookstate"
	conn, _   = dbr.Open("mysql", "username:api@/book_api", nil)
	sess      = conn.NewSession(nil)
)

type server struct{}

type listall []*ListUp

type ListUp struct {
	Id    int    `db:"id"`
	Title string `db:"title"`
}

type Book struct {
	Id     int    `db:"id"`
	Title  string `db:"title"`
	Author string `db:"author"`
	Isbn13 string `db:"isbn13"`
	State  bool   `db:"state"`
	Pic    string `db:"pic"`
}

func arrayBook(pbbook *pb.Book, book *Book) error {
	if pbbook == nil {
		return fmt.Errorf("pbbook is nil")
	}
	*pbbook = pb.Book{
		Id:     strconv.Itoa(book.Id),
		Title:  book.Title,
		Author: book.Author,
		Isbn13: book.Isbn13,
		State:  book.State,
		Pic:    book.Pic,
	}
	return nil
}

//func arraylist(lists *[]pb.Listup, books *ListUp) error {
//	for _, i := range books {
//		lists[i] = pb.Listup{
//			Id:    books[i].Id,
//			Title: books[i].Title,
//		}
//	}
//	return nil
//}

func (s *server) Fetch(ctx context.Context, in *pb.FetchRequest) (*pb.Book, error) {
	var book Book
	pbbook := pb.Book{}
	id := in.Id
	sess.Select("*").From(tablename).Where("id = ?", id).Load(&book)
	if err := arrayBook(&pbbook, &book); err != nil {
		return nil, err
	}
	return &pbbook, nil
}

func (s *server) Add(ctx context.Context, in *pb.Book) (*pb.Empty, error) {
	sess.InsertInto(tablename).Columns("title", "author", "isbn13", "state", "pic").Values(in.Title, in.Author, in.Isbn13, in.State, in.Pic).Exec()
	return &pb.Empty{}, nil
}

func (s *server) ListAll(ctx context.Context, in *pb.Empty) (*pb.ListAllResponse, error) {
	var books listall
	var lists []*pb.Listup
	sess.Select("id, title").From(tablename).Load(&books)

	for _, res := range books {
		lists = append(lists, &pb.Listup{
			Id:    strconv.Itoa(res.Id),
			Title: res.Title,
		})
	}

	return &pb.ListAllResponse{lists}, nil
}

func (s *server) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	attrsMap := map[string]interface{}{"id": in.Book.Id, "title": in.Book.Title, "author": in.Book.Author, "isbn13": in.Book.Isbn13, "state": in.Book.State, "pic": in.Book.Pic}
	sess.Update(tablename).SetMap(attrsMap).Where("id = ?", in.Book.Id).Exec()
	return &pb.UpdateResponse{}, nil
}

func (s *server) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.Empty, error) {
	sess.DeleteFrom(tablename).Where("id = ?", in.Id).Exec()
	return &pb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	defer s.Stop()
}
