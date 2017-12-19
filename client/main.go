package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	pb "github.com/soeyusuke/grpcdemo/proto"
	"google.golang.org/grpc"
)

const (
	port = "1111"
)

var (
	data []string
)

type BookInfo struct {
	id     string
	title  string
	author string
	isbn13 string
	state  string
	pic    string
}

func main() {
	var mode string
	s := BookInfo{}
	flag.StringVar(&mode, "mode", "", "[fetch|add|list|update|delete]を指定")
	flag.StringVar(&s.id, "id", "", "idを指定")
	flag.StringVar(&s.title, "title", "", "titleを指定")
	flag.StringVar(&s.author, "author", "", "authorを指定")
	flag.StringVar(&s.isbn13, "isbn13", "", "isbn13を指定| ture or false")
	flag.StringVar(&s.state, "state", "", "stateを指定")
	flag.StringVar(&s.pic, "pic", "", "picを指定")
	flag.Parse()

	conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer conn.Close()

	c := pb.NewServiceClient(conn)
	ctx := context.Background()

	switch mode {
	case "fetch":
		if s.id == "" {
			log.Fatalf("not set id")
		}
		if err := s.fetch(ctx, &c); err != nil {
			log.Fatal(err)
		}

	case "add":
		if s.title == "" || s.author == "" || s.isbn13 == "" || s.state == "" || s.pic == "" {
			log.Fatalf("no set params are title, author, isbn13, state, pic")
		}
		if err := s.add(ctx, &c); err != nil {
			log.Fatal(err)
		}
		fmt.Println("success add the params")

	case "list":
		if err := s.list(ctx, &c); err != nil {
			log.Fatal(err)
		}

	case "update":
		if s.id == "" || s.title == "" || s.author == "" || s.isbn13 == "" || s.state == "" || s.pic == "" {
			log.Fatalf("no set params are id, title, author, isbn13, state, pic")
		}
		if err := s.update(ctx, &c); err != nil {
			log.Fatal(err)
		}
		fmt.Println("success update the params")

	case "delete":
		if s.id == "" {
			log.Fatalf("not set id")
		}
		if err := s.delete(ctx, &c); err != nil {
			log.Fatal(err)
		}
		fmt.Println("delete success")
	default:
		log.Fatal("Please set the mode [fetch/add/list/update/delete]\nLike this '-mode=add -id=10'")
	}
}

func (s BookInfo) fetch(ctx context.Context, c *pb.ServiceClient) error {
	res, err := (*c).Fetch(ctx, &pb.FetchRequest{
		Id: s.id,
	})
	if err != nil {
		return err
	}
	data := []string{
		res.Id,
		res.Title,
		res.Author,
		res.Isbn13,
		fmt.Sprintf("%v", res.State),
		res.Pic,
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"id", "title", "author", "isbn13", "state", "picture"})
	table.Append(data)
	fmt.Println("Fetch:")
	table.Render()
	return nil
}

func (s BookInfo) add(ctx context.Context, c *pb.ServiceClient) error {
	pbbook := pb.Book{}
	if err := arraydata(&pbbook, s); err != nil {
		return err
	}

	_, err := (*c).Add(ctx, &pbbook)
	if err != nil {
		return err
	}
	return nil
}

func (s BookInfo) list(ctx context.Context, c *pb.ServiceClient) error {
	res, err := (*c).ListAll(ctx, &pb.Empty{})
	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"id", "title"})

	for _, v := range res.Books {
		table.Append([]string{v.Id, v.Title})
	}

	fmt.Println("ListAll:")
	table.Render()
	return nil
}

func (s BookInfo) update(ctx context.Context, c *pb.ServiceClient) error {
	pbbook := pb.Book{}
	if err := arraydata(&pbbook, s); err != nil {
		return err
	}

	_, err := (*c).Update(ctx, &pb.UpdateRequest{&pbbook})
	if err != nil {
		return err
	}
	return nil
}

func (s BookInfo) delete(ctx context.Context, c *pb.ServiceClient) error {
	_, err := (*c).Delete(ctx, &pb.DeleteRequest{
		Id: s.id,
	})
	if err != nil {
		return err
	}
	return nil
}

func arraydata(pbbook *pb.Book, b BookInfo) error {
	*pbbook = pb.Book{
		Id:     b.id,
		Title:  b.title,
		Author: b.author,
		Isbn13: b.isbn13,
		State:  confilmState(b.state),
		Pic:    b.pic,
	}
	return nil
}

func confilmState(j string) bool {
	var judg bool
	if j == "true" {
		judg = true
		return judg
	} else {
		judg = false
		return judg
	}
}
