package main

import (
	"context"
	"fmt"
	"io"
	"log"

	//Paths go to GO ROOT

	"github.com/simplesteph/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Blog client")

	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	//CREATE BLOG
	fmt.Println("Create Blog")
	blog := &blogpb.Blog{
		AuthorId: "Bartosz",
		Title:    "My First blog",
		Content:  "Conetent of the blog",
	}

	createBlog, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Error: %v ", err)
	}

	fmt.Printf("Created the blog: %v\n", createBlog)

	//READ BLOG
	fmt.Println("Read Blog")

	blogID := createBlog.GetBlog().GetId() //can make any string you want
	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlog, err2 := c.ReadBlog(context.Background(), readBlogReq)
	if err2 != nil {
		log.Fatalf("Error: %v ", err2)
	}
	fmt.Printf("Read the blog: %v\n", readBlog)

	//UPDATE BLOG
	fmt.Println("Update Blog")

	blogUpdated := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Bartosz  Zarzecki",
		Title:    "My First blog updated",
		Content:  "Conetent of the blog that is new",
	}

	updateBlogReq := &blogpb.UpdateBlogRequest{Blog: blogUpdated}
	updateBlog, err2 := c.UpdateBlog(context.Background(), updateBlogReq)
	if err2 != nil {
		log.Fatalf("Error: %v ", err2)
	}
	fmt.Printf("Update the blog: %v\n", updateBlog)

	//DELETE BLOG
	fmt.Println("Update Blog")

	deleteBlogReq := &blogpb.DeleteBlogRequest{BlogId: blogID}
	deleteBlog, err2 := c.DeleteBlog(context.Background(), deleteBlogReq)
	if err2 != nil {
		log.Fatalf("Error: %v ", err2)
	}
	fmt.Printf("Deleted the blog: %v\n", deleteBlog)

	//LIST AL THE BLOGS

	fmt.Println("List All Blogs")

	req := &blogpb.ListBlogRequest{}

	res, err := c.ListBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v ", err)
	}

	for {
		msg, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while calling ListBlog RPC: %v ", err)
		}
		resultMessage := msg.GetBlog()
		fmt.Println("Response from Listing \n", resultMessage)
	}

}
