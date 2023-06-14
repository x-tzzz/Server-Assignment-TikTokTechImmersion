package db

import (
	"context"
	"database/sql"
	"time"

	api "xtz-Server-Assignment-TikTokTechImmersion/api"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase() (*Database, error) {
	// Open database
	db, err := sql.Open("sqlite3", "./messages.db")
	if err != nil {
		return nil, err
	}

	// Create messages table if it does not exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS messages (chat TEXT, text TEXT, sender TEXT, send_time INTEGER)")
	if err != nil {
		return nil, err
	}

	// Create index on chat and send_time columns of messages table
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_messages_chat_send_time ON messages(chat, send_time)")
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (d *Database) Pull(ctx context.Context, in *api.PullRequest) (*api.PullResponse, error) {
	rows, err := d.DB.QueryContext(ctx, "SELECT chat, text, sender, send_time FROM messages WHERE chat = ? AND send_time >= ? ORDER BY send_time ASC LIMIT ?",
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

func (d *Database) Send(ctx context.Context, in *api.SendRequest) (*api.SendResponse, error) {
	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	_, err := d.DB.ExecContext(ctx, "INSERT INTO messages (chat, text, sender, send_time) VALUES (?, ?, ?, ?)", in.GetChat(), in.GetText(), in.GetSender(), currentTime)
	if err != nil {
		return nil, err
	}
	return &api.SendResponse{}, nil
}
