package handler

import (
	"context"
	"log"
	"os"
	"time"

	"go.etcd.io/etcd/client/v3"

	pb "github.com/ujjwal-shekhar/load_balancer/services/common/genproto/comms"
	"github.com/ujjwal-shekhar/load_balancer/services/common/utils/constants"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	etcdClient *clientv3.Client
	
	TaskRunnerClient 	pb.TaskRunnerClient
	LoadBalancerClient 	pb.LoadBalancerClient

	Load    			int32
}

func NewClient () *Client {
	conn, err := grpc.NewClient(constants.LB_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to load balancer: %v", err)
	}
	lbClient := pb.NewLoadBalancerClient(conn)

	return &Client{
		LoadBalancerClient: lbClient,
	}
}

func (c *Client) Run() {
	// Lets ask the lb for address of a task runner server
	ctx := context.Background()

	startTime := time.Now()
	resp, err := c.LoadBalancerClient.ProcessClientRequest(ctx, &pb.ClientRequest{Load: c.Load})
	if err != nil {
		log.Fatalf("Failed to get task runner: %v", err)
	}
	log.Println("Got response from lb: ", resp)

	// Connect to the task runner
	conn, err := grpc.NewClient(resp.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to task runner: %v", err)
	}
	taskRunnerClient := pb.NewTaskRunnerClient(conn)

	// Send the task load to the task runner
	log.Printf("Sending task load to task runner: %s at time %v", resp.Address, time.Now().Unix())
	_, err = taskRunnerClient.RunTask(ctx, &pb.ClientRequest{Load: c.Load})

	if err != nil {
		log.Fatalf("Failed to send task load: %v", err)
	} else {
		// Log the end time and calculate turnaround time
		endTime := time.Now()
		turnaroundTime := endTime.Sub(startTime).Seconds()

		// Log the turnaround time to client.log
		file, err := os.OpenFile("client.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open client.log: %v", err)
		}
		defer file.Close()

		logger := log.New(file, "", log.LstdFlags)
		logger.Printf("Turnaround time: %.2f seconds", turnaroundTime)
		log.Printf("Finished task load on task runner: %s at time %v", resp.Address, time.Now().Unix())
	}
}