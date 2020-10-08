package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/grpc/calculator/calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// doServerStreaming(c)

	// doClientStreaming(c)

	// doBidiStreaming(c)
	doErrorUnary(c)
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

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("STARTING to do a COMPUTEAVERAGE CLIENT STREAMINGRPC...")
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream %v", err)
	}

	numbers := []int64{3, 5, 9, 54, 23}

	for _, number := range numbers {
		fmt.Printf("Sending number %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response %v", err)
	}
	fmt.Printf("The average is %v\n", res.GetResult())
}

func doBidiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("STrating to do BIDI Strreaming gRPC")

	stream, err := c.FindMax(context.Background())
	if err != nil {
		log.Fatalf("Error while creating the stream %v\n", err)
		return
	}

	waitc := make(chan struct{})

	go func() {
		numbers := []int32{4, 5, 2, 19, 4, 6, 32}

		for _, number := range numbers {
			fmt.Printf("Sending number %v\n", number)
			stream.Send(&calculatorpb.FindMaxRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving stream %v\n", err)
				break
			}
			maximum := res.GetCurrMax()
			fmt.Printf("Recevied a new maxiumum of...: %v\n", maximum)

		}
		close(waitc)

	}()
	<-waitc

}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("STrating to do SQUAROOT UNARY RPC\n ")
	//correct call

	//Correct call
	doErrorCall(c, 10)
	// Error Call
	doErrorCall(c, -2)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	// fmt.Println("DO ERROR CALL WAS CALLED")

	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			//Actual error
			fmt.Printf("Error Message from server: %v\n", respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Printf("We probably sent a negative number\n")
				return
			}
		} else {
			log.Fatalf("BIG ERROR CALLING SqRoot: %v\n", err)
			return
		}
	}

	fmt.Printf("Result of Square Root of Number: %v : %v\n", n, res.GetNumberRoot())

}
