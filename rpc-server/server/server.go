package server

import (
	"context"
	"database/sql"
	"log"
	"net"
	"time"

	api "xtz-Server-Assignment-TikTokTechImmersion/api"

	"google.golang.org/grpc"
	//"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	api.UnimplementedMessageServiceServer
	DB *sql.DB
}

func NewServer(db *sql.DB) *Server {
	return &Server{DB: db}
}

func (s *Server) Pull(ctx context.Context, in *api.PullRequest) (*api.PullResponse, error) {
	// Assume we have a table "messages" with columns "chat", "text", "sender", "send_time"
	rows, err := s.DB.QueryContext(ctx, "SELECT chat, text, sender, send_time FROM messages WHERE chat = ? AND send_time >= ? ORDER BY send_time ASC LIMIT ?",
		in.GetChat(), in.GetCursor(), in.GetLimit())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*api.Message
	for rows.Next() {
		var msg api.Message
		if err := rows.Scan(&msg.Chat, &msg.Text, &msg.Sender, &msg.SendTime); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}

	return &api.PullResponse{Messages: messages}, nil
}

func (s *Server) Send(ctx context.Context, in *api.SendRequest) (*api.SendResponse, error) {
	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	_, err := s.DB.ExecContext(ctx, "INSERT INTO messages (chat, text, sender, send_time) VALUES (?, ?, ?, ?)", in.GetChat(), in.GetText(), in.GetSender(), currentTime)
	if err != nil {
		return nil, err
	}
	return &api.SendResponse{}, nil
}

func (s *Server) Start(lis net.Listener) error {
	grpcServer := grpc.NewServer()

	// Register our service with the gRPC server
	api.RegisterMessageServiceServer(grpcServer, s)

	// Serve our gRPC server on the listener interface
	return grpcServer.Serve(lis)
}

func main() {
	// Open database
	db, err := sql.Open("sqlite3", "./messages.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS messages (chat TEXT, text TEXT, sender TEXT, send_time INTEGER)")
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	server := NewServer(db)

	// Listen on port 8888
	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("Failed to listen on port 8888: %v", err)
	}

	// Serve our gRPC server on the listener interface
	if err := server.Start(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
