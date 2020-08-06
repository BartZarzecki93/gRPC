package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/simplesteph/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Hello Iam client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	//	doUnary(c)
	//doServerStreaming(c)
	doClientStreaming(c)
}

func doClientStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("Satarting to do a client streaming RPC...")

	request := []*greetpb.LongGreetRequest{

		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bartosz",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Zarzec",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bartolo",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bart",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreat RPC: %v ", err)
	}

	for _, req := range request {

		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving longgreet RPC: %v ", err)
	}
	log.Printf("Response from LongGreet: %v", res)

}
func doServerStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("Satarting to do a streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bartosz",
			FlastName: "Zarzecki",
		},
	}

	res, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v ", err)
	}

	for {
		msg, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while calling GreetManyTimes RPC: %v ", err)
		}

		resultMessage := msg.GetResult()
		fmt.Println("Response from GreatManyTimes ", resultMessage)
	}

}

func doUnary(c greetpb.GreetServiceClient) {

	fmt.Println("Satarting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bartosz",
			FlastName: "Zarzecki",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v ", err)
	}
	log.Printf("Response from Greet: %v", res.Result)

}
