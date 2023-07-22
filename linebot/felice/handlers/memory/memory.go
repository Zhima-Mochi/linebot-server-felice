package memory

import (
	"context"
	"encoding/json"
	"fmt"
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

func (m *Memory) Remember(ctx context.Context, userID string, message openai.ChatCompletionMessage) error {
	return m.cacheHandler.RPush(ctx, userID, message)
}

func (m *Memory) Recall(ctx context.Context, userID string, n int) ([]openai.ChatCompletionMessage, error) {
	values, err := m.cacheHandler.LRange(ctx, userID, -1*int64(n), -1)
	if err != nil {
		return nil, err
	}
	messages := []openai.ChatCompletionMessage{}
	for _, value := range values {
		var message openai.ChatCompletionMessage
		if err := json.Unmarshal([]byte(value), &message); err != nil {
			return nil, fmt.Errorf("failed to unmarshal value for userID '%s': %v", userID, err)
		}

		messages = append(messages, message)
	}
	return messages, nil
}

func (m *Memory) Revoke(ctx context.Context, userID string, n int) ([]openai.ChatCompletionMessage, error) {
	values, err := m.cacheHandler.LRange(ctx, userID, -1*int64(n), -1)
	if err != nil {
		return nil, err
	}
	err = m.cacheHandler.LTrim(ctx, userID, 0, -1*int64(n+1))
	if err != nil {
		return nil, err
	}
	messages := []openai.ChatCompletionMessage{}
	for _, value := range values {
		var message openai.ChatCompletionMessage
		if err := json.Unmarshal([]byte(value), &message); err != nil {
			return nil, fmt.Errorf("failed to unmarshal value for userID '%s': %v", userID, err)
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (m *Memory) Forget(ctx context.Context, userID string, n int) error {
	return m.cacheHandler.LTrim(ctx, userID, int64(n), -1)
}

func (m *Memory) GetSize(ctx context.Context, userID string) (int, error) {
	l, err := m.cacheHandler.LLen(ctx, userID)
	return int(l), err
}
