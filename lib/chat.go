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

func Gippity(question string, info string, maxTokens int) *openai.ChatCompletionStream {
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
	return stream
	// defer stream.Close()
	//
	// fmt.Printf("Stream response: ")
	// for {
	// 	response, err := stream.Recv()
	// 	if errors.Is(err, io.EOF) {
	// 		fmt.Println("\nStream finished")
	// 		return nil
	// 	}
	//
	// 	if err != nil {
	// 		fmt.Printf("\nStream error: %v\n", err)
	// 		return nil
	// 	}
	//
	// 	fmt.Printf(response.Choices[0].Delta.Content)
	// }
}
