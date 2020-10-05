package main

import (
	"context"
	"fmt"
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

	c := calculatorpb.NewSumServiceClient(cc)
	// fmt.Printf("Created clinet %f", c)

	doUnary(c)
}

func doUnary(c calculatorpb.SumServiceClient) {
	fmt.Println("STARTING to do a SUM UNARY RPC...")

	req := &calculatorpb.SumRequest{
		FirstNumber:  5,
		SecondNumber: 40,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling SUM RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.SumResult)
}
