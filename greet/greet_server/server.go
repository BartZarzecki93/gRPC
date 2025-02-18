package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	//Paths go to GO ROOT
	"github.com/simplesteph/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {

	fmt.Printf("Greet with Deadline function was invoked with %v\n", req)

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Printf("Canceled on a client\n")
			return nil, status.Error(codes.Canceled, "the client canceled the request")

		}
		time.Sleep(1 * time.Second)
	}

	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("Greet Everyone function was invoked with with streamign request\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while calling client stream RPC: %v ", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + " !"

		senderr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if senderr != nil {
			log.Fatalf("error while sendig data to  client stream RPC: %v ", err)
			return err
		}

	}

}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with with streamign request\n")
	result := "Hello "
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("error while calling client stream RPC: %v ", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += firstName + " !"

	}
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {

	fmt.Printf("Greet function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {

	fmt.Printf("Greet stream function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManytimesResponse{
			Result: result,
		}
		stream.Send(res)
		//time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	certFile := "ssl/server.crt"
	keyFile := "ssl/server.pem"
	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		log.Fatalf("Failed loading certs: %v ", sslErr)
		return
	}

	opts := grpc.Creds(creds)
	s := grpc.NewServer(opts)
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
