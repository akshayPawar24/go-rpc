package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	ratepb "go-rpc/grpc/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	r := gin.Default()
	r.GET("/rate", func(c *gin.Context) {
		base := c.Query("base")
		target := c.Query("target")
		conn, err := grpc.Dial("localhost:50051",
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		client := ratepb.NewRateServiceClient(conn)

		// Prepare the request
		req := &ratepb.GetRateRequest{
			Base:   base,
			Target: target,
		}

		// Set a timeout for the request
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Call the GetRate method
		resp, err := client.GetRate(ctx, req)
		if err != nil {
			log.Fatalf("could not get rate: %v", err)
		}

		// Print the response
		if resp.Error != "" {
			fmt.Printf("Error: %s\n", resp.Error)
		} else {
			fmt.Printf("Rate: 1 %s = %f %s (updated at %d)\n", resp.Base, resp.Rate, resp.Target, resp.UpdatedAt)
		}
		c.String(200, resp.String())
	})
	r.Run(":8090")
}
