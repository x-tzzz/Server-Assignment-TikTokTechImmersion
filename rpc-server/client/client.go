package client

import (
	"context"
	"log"
	"time"

	api "xtz-Server-Assignment-TikTokTechImmersion/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Service api.MessageServiceClient
}

func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{Service: api.NewMessageServiceClient(conn)}
}

func PullMessages(client api.MessageServiceClient, chat string, cursor int64, limit int32) error {
	req := &api.PullRequest{
		Chat:   chat,
		Cursor: cursor,
		Limit:  limit,
	}

	res, err := client.Pull(context.Background(), req)
	if err != nil {
		return err
	}

	for _, msg := range res.Messages {
		log.Printf("Received message from %s at %s: %s", msg.Sender, time.Unix(msg.SendTime, 0).String(), msg.Text)
	}

	return nil
}

func PushMessage(client api.MessageServiceClient, chat string, text string, sender string) {

	req := &api.SendRequest{
		Chat:   chat,
		Text:   text,
		Sender: sender,
	}

	resp, err := client.Send(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	log.Printf("Response from server: %v", resp)
}

func main() {
	// Connect to chat service
	conn, err := grpc.Dial("localhost:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to chat service: %v", err)
	}
	defer conn.Close()

	// Initialize gRPC client
	client := api.NewMessageServiceClient(conn)

	// Call the PullMessages function
	err = PullMessages(client, "chat1", 0, 10)
	if err != nil {
		log.Fatalf("Failed to pull messages: %v", err)
	}
}
