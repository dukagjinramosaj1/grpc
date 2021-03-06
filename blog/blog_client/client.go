package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/grpc/blog/blogpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}

	defer cc.Close()

	c := blogpb.NewBlogServericeClient(cc)
	// fmt.Printf("Created clinet %f", c)

	//CREATE BLOG
	fmt.Println("Creating the Blog")

	blog := &blogpb.Blog{
		AuthorId: "Dukagjin",
		Title:    "My first Blog",
		Content:  "Content of the first blog",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected Error: %v\n", err)
	}

	fmt.Printf("Blog has been created: %v\n", createBlogRes)
	blogID := createBlogRes.GetBlog().GetId()

	//Read Blog
	fmt.Println("READING THE BLOG")
	_, errRes := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "5f87332e8v4d0653463d539d"})
	if errRes != nil {
		fmt.Printf("Error happened while reading: %v\n", errRes)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error happened while reading: %v\n", readBlogErr)
	}
	fmt.Printf("Blog was read: %v\n", readBlogRes)

	//UPDATE Blog

	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Changed Author",
		Title:    "My first Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions",
	}
	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v\n", updateErr)
	}

	fmt.Printf("Blog was Updated: %v\n", updateRes)

	//Delete Blog

	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if deleteErr != nil {
		fmt.Printf("Error happened while Deleting: %v\n", deleteErr)
	}
	fmt.Printf("Blog was deleted: %v\n", deleteRes)

	//LIST BLOG

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("Error while calling ListBLOG RPC: %v\n", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happned: %v\n", err)
		}
		fmt.Println(res.GetBlog())
	}

}
