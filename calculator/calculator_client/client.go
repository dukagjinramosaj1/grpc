package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/grpc/calculator/calculatorpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("HELLO WORLD")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("coul not conect : %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// fmt.Printf("Created clinet %f", c)

	// doUnary(c)
	doServerStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("STARTING to do a SUM UNARY RPC...")

	req := &calculatorpb.SumRequest{
		FirstNumber:  5,
		SecondNumber: 40,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling SUM RPC: %v", err)
	}
	log.Printf("Response from SUM: %v", res.SumResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("STARTING to do a PRIMEDECOMPOSITION SERVER STREAMINGRPC...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 36,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PRIMEDECOMPOSITION RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("SOMETHING HAPPENED %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}
