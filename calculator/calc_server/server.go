package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/simplesteph/grpc-go-course/calculator/calcpb"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) SquareRoot(ctx context.Context, req *calcpb.SquareRootRequest) (*calcpb.SquareRootResponse, error) {

	fmt.Printf("Received Squreroot %v\n", req)
	firstNumber := req.GetNumber()

	if firstNumber < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("received a negative number: %v ", firstNumber),
		)
	}
	res := &calcpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(firstNumber)),
	}
	return res, nil

	// secondNumber := req.SecondNumber

	// sum := firstNumber + secondNumber
	// res := &calcpb.SumResponse{
	// 	Result: sum,
	// }
	// return res, nil
}

func (*server) FindMaximum(stream calcpb.CalcService_FindMaximumServer) error {

	fmt.Printf("FindMax function was invoked with with streamign request\n")
	initialNumber := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while calling client stream RPC: %v ", err)
			return err
		}

		firstNumber := req.GetNumber()

		if initialNumber < firstNumber {
			initialNumber = firstNumber
			senderr := stream.Send(&calcpb.FindMaximumResponse{
				Maximum: firstNumber,
			})
			if senderr != nil {
				log.Fatalf("error while sendig data to  client stream RPC: %v ", err)
				return err
			}
		}

	}

}

func (*server) ComputeAverage(stream calcpb.CalcService_ComputeAverageServer) error {

	fmt.Printf("ComputeAve function was invoked with with streamign request\n")
	result := float64(0)
	count := float64(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calcpb.ComputeAverageResponse{
				Average: result / count,
			})
		}
		if err != nil {
			log.Fatalf("error while calling client stream RPC: %v ", err)
		}

		result = result + float64(req.GetNumber())
		count++

	}

}

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

	reflection.Register(s)
	
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
