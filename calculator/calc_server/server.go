package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/simplesteph/grpc-go-course/calculator/calcpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calcpb.SumRequest) (*calcpb.SumResponse, error) {

	fmt.Printf("Received sum %v\n", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber

	sum := firstNumber + secondNumber
	res := &calcpb.SumResponse{
		Result: sum,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calcpb.PrimeNumberDecompositionRequest, stream calcpb.CalcService_PrimeNumberDecompositionServer) error {

	fmt.Printf("Received numbers from stream  %v\n", req)
	number := req.GetNumber()
	k := int64(2)

	for number > 1 {
		if number%k == 0 {
			stream.Send(&calcpb.PrimeNumberDecompositionResponse{
				FractorNumber: k,
			})
			number = number / k
		} else {
			k++
			fmt.Printf("Divisor has increased\n")
		}
	}

	return nil
}

func main() {
	fmt.Println("Calc Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calcpb.RegisterCalcServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
