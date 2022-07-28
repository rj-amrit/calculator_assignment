package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	calculator "github.com/rj-amrit/calculator_proto/calculator"
	grpc "google.golang.org/grpc"
)

type server struct {
	calculator.UnimplementedCalculatorServer
}

func (*server) Sum(ctx context.Context, req *calculator.TwoNumRequest) (resp *calculator.SumResponse, err error) {
	fmt.Println("Sum function is invoked...")
	num1 := req.GetNum1()
	num2 := req.GetNum2()
	num3 := num1 + num2
	resp = &calculator.SumResponse{
		Sum: num3,
	}
	return resp, nil
}

func (*server) PrimeNumbers(req *calculator.NumRequest, resp calculator.Calculator_PrimeNumbersServer) error {
	fmt.Println("PrimeNumbers function is invoked...")
	num := req.GetNum()
	var i int64
	for i = 2; i < num; i++ {
		var j int64
		flag := false
		for j = 2; j < i; j++ {
			if i%j == 0 {
				flag = true
				break
			}
		}
		if !flag {
			res := calculator.AllPrimesResponse{
				Num: i,
			}
			resp.Send(&res)
		}
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) ComputeAverage(stream calculator.Calculator_ComputeAverageServer) error {
	fmt.Println("ComputeAverage function is invoked...")
	var sum int64
	var count int64
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculator.AverageResponse{
				Num: float32(sum) / float32(count),
			})
		}
		if err != nil {
			log.Fatal(err)
		}
		sum += msg.GetNum()
		count++
	}
}

func (*server) FindMaxNumber(stream calculator.Calculator_FindMaxNumberServer) error {
	fmt.Println("FindMaxNumber function is invoked...")
	var maxNum int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatal(err)
			return err
		}
		num := req.GetNum()
		if num > maxNum {
			maxNum = num
		}
		stream.Send(&calculator.MaxNumResponse{
			Num: maxNum,
		})
	}
}
func main() {
	fmt.Println("vim-go")

	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Failed to Listen: %v", err)
	}

	s := grpc.NewServer()
	calculator.RegisterCalculatorServer(s, &server{})

	if err = s.Serve(listen); err != nil {
		log.Fatal("failed to serve : %v", err)
	}
}
