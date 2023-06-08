package lib

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type Chat struct {
	text string
	*openai.ChatCompletionStream
}

func (chat *Chat) Text() (string, error) {
	response, err := chat.Recv()
	if err != nil {
		return "", err
	}
	chat.text += response.Choices[0].Delta.Content
	return chat.text, err
}

func Gippity(question string, info string, maxTokens int) *Chat {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	c := openai.NewClient(os.Getenv("OPENAI_KEY"))

	knowledgeEmbeddings := CreateEmbeddings(c, info, maxTokens)
	questionEmbeddings := CreateEmbeddings(c, question, maxTokens)
	topInfo := TopMatches(questionEmbeddings, knowledgeEmbeddings, maxTokens)

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
	stream, err := c.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return nil
	}
	return &Chat{"", stream}
}
