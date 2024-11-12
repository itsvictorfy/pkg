package chatgpt

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type ChatGPTClient struct {
	*openai.Client `json:"-"` // Exclude from JSON
}

// SendChatGptRequest sends a request to ChatGPT with the given prompt and message
func (c *ChatGPTClient) SendChatGptRequest(prompt, msg string) (string, error) {
	req := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: prompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: msg,
		},
	}
	resp, err := c.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4o,
			Messages: req,
		},
	)
	if err != nil {
		return "", fmt.Errorf("sending request to ChatGPT failed: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}
