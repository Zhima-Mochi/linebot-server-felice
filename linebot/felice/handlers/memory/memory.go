package memory

import (
	"context"
	"linebot-server-felice/handlers/cache"

	"github.com/sashabaranov/go-openai"
)

type Memory struct {
	cacheHandler *cache.CacheHandler
}

func NewMemory() *Memory {
	return &Memory{
		cacheHandler: cache.NewCacheHandler(),
	}
}

func (m *Memory) Remember(ctx context.Context, userID string, message openai.ChatCompletionMessage) {
	m.cacheHandler.LPush(ctx, userID, message)
}

func (m *Memory) Recall(ctx context.Context, userID string, n int) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{}
	m.cacheHandler.LRange(ctx, userID, -1*int64(n), -1, &messages)
	return messages
}

func (m *Memory) Revoke(ctx context.Context, userID string, n int) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{}
	m.cacheHandler.LRange(ctx, userID, -1*int64(n), -1, &messages)
	m.cacheHandler.LTrim(ctx, userID, 0, -1*int64(n-1))
	return messages
}
