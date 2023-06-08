package lib

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

var client *openai.Client

type Chat struct {
	messages  []string
	knowledge map[string][]float32
	*openai.ChatCompletionStream
}

func (chat *Chat) Messages() ([]string, error) {
	response, err := chat.Recv()
	if err != nil {
		return chat.messages, err
	}
	chat.messages[len(chat.messages)-1] += response.Choices[0].Delta.Content
	return chat.messages, err
}

func (chat *Chat) SetKnowledge(knowledge string) {
	chat.knowledge = CreateEmbeddings(client, knowledge)
}

func (chat *Chat) Ask(question string, maxTokens int) {
	chat.messages = append(chat.messages, "")
	topInfo := TopMatches(CreateEmbeddings(client, question), chat.knowledge, maxTokens)

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 200,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: question,
			},
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: fmt.Sprintf("Please use the provided context to answer any subsequent questions.: %s", strings.Join(topInfo, "")),
			},
		},
		Stream: true,
	}
	stream, err := client.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		fmt.Println("Error crating chat completion stream: ", err)
	}
	chat.ChatCompletionStream = stream
}

func Init(key string) {
	if client == nil {
		client = openai.NewClient(key)
	}
}

func NewChat(knowledge string) *Chat {
	return &Chat{[]string{}, CreateEmbeddings(client, knowledge), nil}
}
