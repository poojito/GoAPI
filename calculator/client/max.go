package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/Clement-Jean/grpc-go-course/calculator/proto"
)

func doMax(c pb.CalculatorServiceClient) {
	log.Println("doMax was invoked")
	stream, err := c.Max(context.Background())

	if err != nil {
		log.Fatalf("error while opening stream: %v\n", err)
	}
	waitc := make(chan struct{})

	go func() {
		numbers := []int32{4, 5, 2, 19, 4, 6, 32}
		for _, number := range numbers {
			log.Println("Sending number: %d\n", number)
			stream.Send(&pb.MaxRequest{
				Number: number,
			})
			time.Sleep(1 * time.Second)
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
				log.Printf("Problem while reading server stream: %v", err)
				break
			}
			log.Printf("Received a new maximum: %d\n", res.Result)
		}
		close(waitc)
	}()
	<-waitc
}
