package handler

import (
	"context"
	"log"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
	_, err = taskRunnerClient.RunTask(ctx, &pb.ClientRequest{Load: c.Load})

	if err != nil {
		log.Fatalf("Failed to send task load: %v", err)
	}

	log.Printf("Task load sent to task runner: %s", resp.Address)
}