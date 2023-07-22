package memory

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemory(t *testing.T) {
	// Start a Redis service
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer client.Close()

	// Ensure Redis service is up
	_, err := client.Ping(context.Background()).Result()
	require.NoError(t, err)
	// clear all keys
	err = client.FlushAll(context.Background()).Err()
	require.NoError(t, err)

	message1 := openai.ChatCompletionMessage{
		Content: "message1",
	}
	message2 := openai.ChatCompletionMessage{
		Content: "message2",
	}
	message3 := openai.ChatCompletionMessage{
		Content: "message3",
	}

	// Create a new Memory instance
	m := NewMemory()

	// Test Remember method
	err = m.Remember(context.Background(), "user1", message1)
	require.NoError(t, err)
	err = m.Remember(context.Background(), "user1", message2)
	require.NoError(t, err)
	err = m.Remember(context.Background(), "user1", message3)
	require.NoError(t, err)

	// Test GetSize method
	size, err := m.GetSize(context.Background(), "user1")
	require.NoError(t, err)
	assert.Equal(t, 3, size, fmt.Sprintf("Expected size to be 1, got %d", size))

	// Test Recall method
	messages, err := m.Recall(context.Background(), "user1", 3)
	require.NoError(t, err)
	assert.Equal(t, []openai.ChatCompletionMessage{message1, message2, message3}, messages, fmt.Sprintf("Expected messages to be %v, got %v", []openai.ChatCompletionMessage{message1, message2, message3}, messages))

	// Test Revoke method
	messages, err = m.Revoke(context.Background(), "user1", 1)
	require.NoError(t, err)
	assert.Equal(t, []openai.ChatCompletionMessage{message3}, messages, fmt.Sprintf("Expected messages to be %v, got %v", []openai.ChatCompletionMessage{message3}, messages))
	messages, err = m.Recall(context.Background(), "user1", 3)
	require.NoError(t, err)
	assert.Equal(t, []openai.ChatCompletionMessage{message1, message2}, messages, fmt.Sprintf("Expected messages to be %v, got %v", []openai.ChatCompletionMessage{message1, message2}, messages))

	// Test Forget method
	err = m.Forget(context.Background(), "user1", 1)
	require.NoError(t, err)
	messages, err = m.Recall(context.Background(), "user1", 3)
	require.NoError(t, err)
	assert.Equal(t, []openai.ChatCompletionMessage{message2}, messages, fmt.Sprintf("Expected messages to be %v, got %v", []openai.ChatCompletionMessage{message2}, messages))
}
