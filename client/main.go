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
	id     string
	title  string
	author string
	isbn13 string
	state  string
	pic    string
	data   []string
)

func main() {
	var mode string
	flag.StringVar(&mode, "mode", "", "[fetch|add|list|update|delete]を指定")
	flag.StringVar(&id, "id", "", "idを指定")
	flag.StringVar(&title, "title", "", "titleを指定")
	flag.StringVar(&author, "author", "", "authorを指定")
	flag.StringVar(&isbn13, "isbn13", "", "isbn13を指定| ture or false")
	flag.StringVar(&state, "state", "", "stateを指定")
	flag.StringVar(&pic, "pic", "", "picを指定")
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
		if id == "" {
			log.Fatalf("not set id")
		}
		if err := fetch(ctx, &c, id); err != nil {
			log.Fatal(err)
		}

	case "add":
		if title == "" || author == "" || isbn13 == "" || state == "" || pic == "" {
			log.Fatalf("no set params are title, author, isbn13, state, pic")
		} else {
			AppendData()
		}
		if err := add(ctx, &c, data); err != nil {
			log.Fatal(err)
		}
		fmt.Println("success add the params")

	case "list":
		if err := list(ctx, &c); err != nil {
			log.Fatal(err)
		}

	case "update":
		if id == "" || title == "" || author == "" || isbn13 == "" || state == "" || pic == "" {
			log.Fatalf("no set params are id, title, author, isbn13, state, pic")
		} else {
			AppendData()
		}
		if err := update(ctx, &c, data); err != nil {
			log.Fatal(err)
		}
		fmt.Println("success update the params")

	case "delete":
		if id == "" {
			log.Fatalf("not set id")
		}
		if err := delete(ctx, &c, id); err != nil {
			log.Fatal(err)
		}
		fmt.Println("delete success")
	default:
		log.Fatal("Please set the mode [fetch/add/list/update/delete]\nLike this '-mode=add -id=10'")
	}
}

func fetch(ctx context.Context, c *pb.ServiceClient, fid string) error {
	res, err := (*c).Fetch(ctx, &pb.FetchRequest{
		Id: fid,
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

func add(ctx context.Context, c *pb.ServiceClient, data []string) error {
	pbbook := pb.Book{}
	if err := arraydata(&pbbook, data); err != nil {
		return err
	}

	_, err := (*c).Add(ctx, &pbbook)
	if err != nil {
		return err
	}
	return nil
}

func list(ctx context.Context, c *pb.ServiceClient) error {
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

func update(ctx context.Context, c *pb.ServiceClient, data []string) error {
	pbbook := pb.Book{}
	if err := arraydata(&pbbook, data); err != nil {
		return err
	}

	_, err := (*c).Update(ctx, &pb.UpdateRequest{&pbbook})
	if err != nil {
		return err
	}
	return nil
}

func delete(ctx context.Context, c *pb.ServiceClient, did string) error {
	_, err := (*c).Delete(ctx, &pb.DeleteRequest{
		Id: did,
	})
	if err != nil {
		return err
	}
	return nil
}

func arraydata(pbbook *pb.Book, data []string) error {
	var judg bool
	if data == nil {
		return fmt.Errorf("data's contents is nothing")
	}
	if data[4] == "true" {
		judg = true
	} else {
		judg = false
	}

	*pbbook = pb.Book{
		Id:     data[0],
		Title:  data[1],
		Author: data[2],
		Isbn13: data[3],
		State:  judg,
		Pic:    data[5],
	}
	return nil
}

func AppendData() {
	data = append(data, id)
	data = append(data, title)
	data = append(data, author)
	data = append(data, isbn13)
	data = append(data, state)
	data = append(data, pic)
}
