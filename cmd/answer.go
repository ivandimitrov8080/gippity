/*
Copyright © 2023 Ivan Dimitrov ivan@idimitrov.dev
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/ivandimitrov8080/gippity/lib"
	"github.com/joho/godotenv"

	"github.com/gosuri/uilive"
	"github.com/spf13/cobra"
)

var answerCmd = &cobra.Command{
	Use:   "answer",
	Short: "Answers the provided question using the provided knowledge-base.",
	Long: `
  Ask a model a question with a provided knowledge-base.
  `,
	Run: func(cmd *cobra.Command, args []string) {
		question, _ := cmd.Flags().GetString("question")
		knowledgeBase, _ := cmd.Flags().GetString("knowledge")
		tokenLimit, _ := cmd.Flags().GetInt("tokens")
		godotenv.Load()
		lib.Init(os.Getenv("OPENAI_KEY"))
		chat := lib.NewChat(knowledgeBase)
		chat.Ask(question, tokenLimit)
		writer := uilive.New()
		writer.Start()
		for {
			msg, err := chat.Messages()
			if err != nil {
				writer.Stop()
				return
			}
			fmt.Fprintf(writer, "%s\n", msg[len(msg)-1])
		}
	},
}

func init() {
	rootCmd.AddCommand(answerCmd)
	answerCmd.Flags().StringP("question", "q", "", "The question you'd like to know the answer of.")
	answerCmd.Flags().StringP("knowledge", "k", "", "The custom information to take into account.")
	answerCmd.Flags().IntP("tokens", "t", 3000, "Token limit.")
}
