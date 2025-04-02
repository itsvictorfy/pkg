package chatgpt

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type ChatGPTClient struct {
	*openai.Client `json:"-"` // Exclude from JSON
}

func InitClient(apiKey string) (*ChatGPTClient, error) {
	client := openai.NewClient(apiKey) // Create the OpenAI client
	return &ChatGPTClient{Client: client}, nil
}

// SendChatGptRequest sends a request to ChatGPT with the given prompt and message
func (c *ChatGPTClient) SendChatGptRequest(userPrompt, msg, reasoningEffort string, jsonSchema *openai.ChatCompletionResponseFormat) (string, error) {
	req := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: userPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: msg,
		},
	}
	resp, err := c.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:           openai.O3Mini,
			Messages:        req,
			ReasoningEffort: reasoningEffort,
			ResponseFormat:  jsonSchema,
		},
	)
	if err != nil {
		return "", fmt.Errorf("sending request to ChatGPT failed: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}
func (c *ChatGPTClient) SendChatGptRequestWithHistory(userPrompt string, history []openai.ChatCompletionMessage) (string, []openai.ChatCompletionMessage, error) {
	// Append the user's prompt to the conversation history
	history = append(history, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userPrompt,
	})

	// Send the request to the ChatGPT API
	resp, err := c.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4, // Correct model name
			Messages: history,
		},
	)
	if err != nil {
		return "", nil, fmt.Errorf("sending request to ChatGPT failed: %w", err)
	}

	// Ensure the response has at least one choice
	if len(resp.Choices) == 0 {
		return "", nil, fmt.Errorf("no response received from ChatGPT")
	}

	// Append the assistant's response to the history
	history = append(history, resp.Choices[0].Message)

	// Return the assistant's response and updated history
	return resp.Choices[0].Message.Content, history, nil
}

func InitChatGptHistory(globalPrompt string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem, // "system" role for global instructions
			Content: globalPrompt,
		},
	}
}
func CreateChatResponseFormat(schema json.Marshaler, strict bool) *openai.ChatCompletionResponseFormat {
	return &openai.ChatCompletionResponseFormat{
		Type: "json_object", // Assuming this is the correct type
		JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
			Name:        "my_schema",
			Description: "Global JSON Schema",
			Schema:      schema,
			Strict:      strict,
		},
	}
}
