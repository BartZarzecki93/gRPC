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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	//doBiDi(c)
	doUnaryWithDealine(c, 5*time.Second)
	doUnaryWithDealine(c, 1*time.Second)
}

func doUnaryWithDealine(c greetpb.GreetServiceClient, timeout time.Duration) {

	fmt.Println("Satarting to do a Unary with Dealine RPC...")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bartosz",
			FlastName: "Zarzecki",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			//actual error from gRPC    user error
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.DeadlineExceeded {
				fmt.Println("We took too long! Dealine was excided")
			}
		} else {
			log.Fatalf("error while calling Func with deadline RPC: %v ", err)
		}
		return
	}
	log.Printf("Response from Greet with dealine: %v", res.Result)

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
