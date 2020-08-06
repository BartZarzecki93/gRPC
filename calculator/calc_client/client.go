package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/simplesteph/grpc-go-course/calculator/calcpb"
	"google.golang.org/grpc"
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
	doServerStreaming(c)
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
