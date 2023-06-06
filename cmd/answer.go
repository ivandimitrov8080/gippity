/*
Copyright © 2023 Ivan Dimitrov ivan@idimitrov.dev
*/
package cmd

import (
	"github.com/ivandimitrov8080/gippity/lib"

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
		lib.Gippity(question, knowledgeBase, tokenLimit)
	},
}

func init() {
	rootCmd.AddCommand(answerCmd)
	answerCmd.Flags().StringP("question", "q", "", "The question you'd like to know the answer of.")
	answerCmd.Flags().StringP("knowledge", "k", "", "The custom information to take into account.")
	answerCmd.Flags().IntP("tokens", "t", 3000, "Token limit.")
}
