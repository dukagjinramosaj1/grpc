package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/grpc/greet/greetpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Localhost web server. Running on port 50051...")

	tls := true
	opts := grpc.WithInsecure()
	if tls {
		certFile := "../../ssl/ca.crt" //certificate authority trust certificate

		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
			return
		}

		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("coul not conect : %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created clinet %f", c)

	doUnary(c)

	// doServerStreaming(c)

	// doClientStreaming(c)

	// doBidiStreaming(c)

	// doUnaryWithDeadline(c, 5*time.Second) //should complete

	// doUnaryWithDeadline(c, 1*time.Second) // should timeout
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("STARTING to do a UNARY RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Dukagjin",
			LastName:  "Ramosaj",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("STrating to do Server Strreaming RPC")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Dukagjin",
			LastName:  "Ramosaj",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Erro while calling Greetmanytimes RPC: %v", err)
	}

	for {

		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Errro while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())

	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("STrating to do Client Strreaming RPC")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Dukagjin",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Duke",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Rob",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bob",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error whiel calling LongGreet %v", err)
	}

	//We iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending requests: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error whiel receiving response from LongGreet %v", err)
	}

	fmt.Printf("LongGreet response : %v\n", resp)
}

func doBidiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("STrating to do BIDI Strreaming RPC")

	//We create a stream by invoking the clinet

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream %v\n", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Dukagjin",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Duke",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Rob",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bob",
			},
		},
	}

	waitc := make(chan struct{})
	// we send bunc of messages to the client (go routine)
	go func() {
		//function to send a bunch of messages

		for _, req := range requests {
			fmt.Printf("SEnding message: %v \n ", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		//function to Receive  a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error while receiving stream %v\n", err)
				break
			}
			fmt.Printf("Received: %v\n ", res.GetResult())

		}
		close(waitc)
	}()
	//block until everything is done
	<-waitc

}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("STARTING to do a UNARYWITHDEADLINE RPC...")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Dukagjin",
			LastName:  "Ramosaj",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("TIMEOUT WAS HIT, DEADLINE WAS EXCEEDED")
			} else {

				fmt.Printf("Unexpected error: %v", statusErr)
			}
		} else {

			log.Fatalf("Error while calling GreetWithDeadline RPC: %v", err)
		}
		return
	}
	log.Printf("Response from Greet: %v", res.Result)
}
