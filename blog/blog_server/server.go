package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/simplesteph/grpc-go-course/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

var collection *mongo.Collection

type server struct{}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Contet   string             `bson:"content"`
	Title    string             `bson:"title"`
}

func main() {
	//if we crash we got file name and lane number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//connect to mongodb locally
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	fmt.Println("MongoDb Connected")
	//create db and collection
	collection = client.Database("mydb").Collection("blog")

	//start the server
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})
	fmt.Println("Blog Service Started")

	go func() {

		fmt.Println("Starting Server ...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

	}()

	//wait for controll C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	//Block until signal got received
	<-ch

	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Stopping the listener")
	lis.Close()
	fmt.Println("Stopping the mongodb")
	client.Disconnect(context.TODO())
	fmt.Println("End of program")
}
