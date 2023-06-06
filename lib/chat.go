package lib

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func Gippity(question string, info string, maxTokens int) {
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
				Content: fmt.Sprintf("Please use the provided context with regards to the Java library DataPipeline from NorthConcepts to answer any question you may be asked: {'context': {%s}}", strings.Join(topInfo, "")),
			},
		},
		Stream: true,
	}
	stream, err := c.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	fmt.Printf("Stream response: ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		fmt.Printf(response.Choices[0].Delta.Content)
	}
}