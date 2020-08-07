package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	//Paths go to GO ROOT
	"github.com/simplesteph/grpc-go-course/calculator/calcpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {

	fmt.Println("Hello I am client calculator")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := calcpb.NewCalcServiceClient(cc)

	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDi(c)
	doErrorUnary(c, 10)
	doErrorUnary(c, -10)
}

func doErrorUnary(c calcpb.CalcServiceClient, n int32) {

	fmt.Println("Satarting to do a Unary SquarRoot RPC...")
	req := &calcpb.SquareRootRequest{
		Number: n,
	}
	res, err := c.SquareRoot(context.Background(), req)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			//actual error from gRPC    user error
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We sent negative number")
			}
		} else {
			log.Fatalf("error while calling SR RPC: %v ", err)
		}

	}
	log.Printf("Response from SR: %v", res.GetNumberRoot())

}

func doBiDi(c calcpb.CalcServiceClient) {

	fmt.Println("Satarting to do bidi a client streaming RPC...")

	// instead
	/*
		numbers := []int32{1,5,2,7}

		for _, req := range numbers{

			fmt.Printf("Sendign Message: %v\n", req)
				stream.Send(&calcpb.FindMaximumRequest{
				Number: req,
			},)
				time.Sleep(1000 * time.Millisecond)
		}
	*/

	request := []*calcpb.FindMaximumRequest{

		&calcpb.FindMaximumRequest{
			Number: 1,
		},
		&calcpb.FindMaximumRequest{
			Number: 5,
		},
		&calcpb.FindMaximumRequest{
			Number: 2,
		},
		&calcpb.FindMaximumRequest{
			Number: 4,
		},
	}

	stream, err := c.FindMaximum(context.Background())
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
			fmt.Printf("Receiving Message - new maximum: %v\n", res.GetMaximum())
		}
		close(waitc)
	}()

	//block until everything is done
	<-waitc

}

func doClientStreaming(c calcpb.CalcServiceClient) {

	fmt.Println("Satarting to do a client streaming RPC...")

	request := []*calcpb.ComputeAverageRequest{

		&calcpb.ComputeAverageRequest{
			Number: 1,
		},
		&calcpb.ComputeAverageRequest{
			Number: 2,
		},
		&calcpb.ComputeAverageRequest{
			Number: 3,
		},
		&calcpb.ComputeAverageRequest{
			Number: 4,
		},
	}

	stream, err := c.ComputeAverage(context.Background())
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

func doServerStreaming(c calcpb.CalcServiceClient) {

	fmt.Println("Satarting to do a streaming RPC...")

	req := &calcpb.PrimeNumberDecompositionRequest{
		Number: 230,
	}

	res, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC: %v ", err)
	}

	for {
		msg, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while calling GreetManyTimes RPC: %v ", err)
		}
		resultMessage := msg.GetFractorNumber()
		fmt.Println("Response from Factoring ", resultMessage)
	}

}
func doUnary(c calcpb.CalcServiceClient) {

	fmt.Println("Satarting to do a Unary Sum RPC...")
	req := &calcpb.SumRequest{

		FirstNumber:  10,
		SecondNumber: 3,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v ", err)
	}
	log.Printf("Response from Sum: %v", res.Result)

}
