package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	api "xtz-Server-Assignment-TikTokTechImmersion/api"
)

func TestDatabase(t *testing.T) {
	// 初始化数据库
	db, err := NewDatabase()
	assert.NoError(t, err)
	defer db.DB.Close()

	// Clear database before running tests
	_, err = db.DB.ExecContext(context.Background(), "DELETE FROM messages")
	assert.NoError(t, err)

	// 测试 Send 函数
	sendReq := &api.SendRequest{
		Chat:   "chat1",
		Text:   "hello",
		Sender: "user1",
	}

	_, err = db.Send(context.Background(), sendReq)
	assert.NoError(t, err)

	// 测试 Pull 函数
	pullReq := &api.PullRequest{
		Chat:   "chat1",
		Limit:  10,
		Cursor: 0,
	}

	resp, err := db.Pull(context.Background(), pullReq)
	assert.NoError(t, err)

	// 检查返回的消息
	assert.NotNil(t, resp.Messages)
	assert.Len(t, resp.Messages, 1)
	assert.Equal(t, sendReq.Chat, resp.Messages[0].Chat)
	assert.Equal(t, sendReq.Text, resp.Messages[0].Text)
	assert.Equal(t, sendReq.Sender, resp.Messages[0].Sender)
}
