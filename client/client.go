package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	calculator "github.com/rj-amrit/calculator_proto/calculator"
	grpc "google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculator.NewCalculatorClient(cc)

	// Sum(c)
	// PrimeNumbers(c)
	// ComputeAverage(c)
	FindMaxNumber(c)
}

func Sum(c calculator.CalculatorClient) {
	fmt.Println("Starting Sum service...")
	req := calculator.TwoNumRequest{
		Num1: 10,
		Num2: 20,
	}
	resp, err := c.Sum(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while calling Sum grpc unary call: %v", err)
	}
	log.Printf("Response from Sum grpc unary call: %v", resp.Sum)
}

func PrimeNumbers(c calculator.CalculatorClient) {
	fmt.Println("Starting Prime number service...")
	req := calculator.NumRequest{
		Num: 15,
	}
	resStream, err := c.PrimeNumbers(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumbers server-side streaming grpc: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while receving server stream: %v", err)
		}
		fmt.Println("Response From PrimeNumbers Server: ", msg.Num)
	}

}

func ComputeAverage(c calculator.CalculatorClient) {
	fmt.Println("Starting ComputeAverage service...")
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error occured while performing client-side streaming : %v", err)
	}
	req := []*calculator.NumRequest{
		{Num: 5},
		{Num: 10},
		{Num: 15},
		{Num: 20},
		{Num: 25},
		{Num: 30},
	}
	for _, val := range req {
		fmt.Println("\nSending number...: ", val)
		stream.Send(val)
		time.Sleep(1 * time.Second)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from server : %v", err)
	}
	fmt.Println("\n****Response From Server : ", resp.Num)

}

func FindMaxNumber(c calculator.CalculatorClient) {
	fmt.Println("Starting Maxnumber service")
	req := []*calculator.NumRequest{
		{Num: 1},
		{Num: 3},
		{Num: 7},
		{Num: 9},
		{Num: 2},
		{Num: 5},
		{Num: 22},
		{Num: 15},
		{Num: 21},
		{Num: 19},
	}
	stream, err := c.FindMaxNumber(context.Background())
	if err != nil {
		log.Fatalf("error occured while performing client side streaming : %v", err)
	}

	waitchan := make(chan int32)

	go func(req []*calculator.NumRequest) {
		for _, val := range req {
			fmt.Println("\nSending number...: ", val.Num)
			err := stream.Send(val)
			if err != nil {
				log.Fatalf("error while sending request to FindMaxNumber service : %v", err)
			}
			time.Sleep(1000 * time.Millisecond)

		}
		stream.CloseSend()
	}(req)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(waitchan)
				return
			}
			if err != nil {
				log.Fatalf("Error receiving response from server : %v", err)
			}
			fmt.Printf("\nResponse From Server : %v", resp.Num)
		}
	}()
	<-waitchan
}
