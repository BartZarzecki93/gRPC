package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	//Paths go to GO ROOT
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
	//doClientStreaming(c)
	doBiDi(c)
}

func doBiDi(c greetpb.GreetServiceClient) {
	fmt.Println("Satarting to do a client bidi streaming RPC...")

	request := []*greetpb.GreetEveryoneRequest{

		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bartosz",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Zarzec",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bartolo",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bart",
			},
		},
	}

	//creating a stream by invoking the client

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while calling GreetEveryone stream RPC: %v ", err)
		return
	}

	waitc := make(chan struct{})
	//sending a bunch of messages to the client (go routine)
	go func() {
		for _, req := range request {
			fmt.Printf("Sendign Message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	//we receive a bunch of messages from the client

	go func() {

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving GreetEveryone stream RPC: %v ", err)
				break
			}
			fmt.Printf("Receiving Message: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	//block until everything is done
	<-waitc
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
