package test

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"sync"
	"testing"

	api "xtz-Server-Assignment-TikTokTechImmersion/api"
	server "xtz-Server-Assignment-TikTokTechImmersion/rpc-server/server"

	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const benchmarkBufSize = 1024 * 1024

var benchmarkLis *bufconn.Listener

func init() {
	// Initialize the in-memory buffer connection
	benchmarkLis = bufconn.Listen(benchmarkBufSize)
}

func benchmarkBufDialer(context.Context, string) (net.Conn, error) {
	return benchmarkLis.Dial()
}

func BenchmarkConcurrentUsers(b *testing.B) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(b, err)
	defer db.Close()

	// Create messages table in the in-memory database
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS messages (chat TEXT, text TEXT, sender TEXT, send_time INTEGER)")
	require.NoError(b, err)

	// Create and start server
	srv := server.NewServer(db)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := srv.Start(benchmarkLis)
		require.NoError(b, err)
	}()

	wg.Wait() // wait for server to start

	// Connect to the server
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(benchmarkBufDialer), grpc.WithInsecure())
	require.NoError(b, err)
	defer conn.Close()

	// Create client
	client := api.NewMessageServiceClient(conn)

	b.SetParallelism(20) // Set the number of parallel benchmarks
	b.RunParallel(func(pb *testing.PB) {
		// Each goroutine will send and pull messages until the benchmark is done
		for pb.Next() {
			for i := 0; i < 100; i++ { // Each user sends and pulls 100 messages
				sendReq := &api.SendRequest{
					Chat:   fmt.Sprintf("chat%d", i),
					Text:   "hello",
					Sender: "user1",
				}
				_, err := client.Send(context.Background(), sendReq)
				require.NoError(b, err)

				pullReq := &api.PullRequest{
					Chat:   fmt.Sprintf("chat%d", i),
					Cursor: 0,
					Limit:  10,
				}
				_, err = client.Pull(context.Background(), pullReq)
				require.NoError(b, err)
			}
		}
	})
}
