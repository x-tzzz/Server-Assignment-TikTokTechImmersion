package test

import (
	"context"
	"database/sql"
	"net"
	"testing"

	api "xtz-Server-Assignment-TikTokTechImmersion/api"
	//"xtz-Server-Assignment-TikTokTechImmersion/rpc-server/client"
	server "xtz-Server-Assignment-TikTokTechImmersion/rpc-server/server"

	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	// Initialize the in-memory buffer connection
	lis = bufconn.Listen(bufSize)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestSendAndPull(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS messages (chat TEXT, text TEXT, sender TEXT, send_time INTEGER)")
	require.NoError(t, err)

	// Create and start server
	srv := server.NewServer(db)
	go func() {
		require.NoError(t, srv.Start(lis))
	}()

	// Connect to the server
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	// Create client
	client := api.NewMessageServiceClient(conn)

	// Test send message
	sendReq := &api.SendRequest{
		Chat:   "chat1",
		Text:   "hello",
		Sender: "user1",
	}
	_, err = client.Send(context.Background(), sendReq)
	require.NoError(t, err)

	// Test pull messages
	pullReq := &api.PullRequest{
		Chat:   "chat1",
		Cursor: 0,
		Limit:  10,
	}
	pullRes, err := client.Pull(context.Background(), pullReq)
	require.NoError(t, err)
	require.Len(t, pullRes.Messages, 1)
	require.Equal(t, "hello", pullRes.Messages[0].Text)
	require.Equal(t, "user1", pullRes.Messages[0].Sender)

	// Test pull from non-existent chat
	pullReq = &api.PullRequest{
		Chat:   "chat2",
		Cursor: 0,
		Limit:  10,
	}
	pullRes, err = client.Pull(context.Background(), pullReq)
	require.NoError(t, err)
	require.Len(t, pullRes.Messages, 0)
}
