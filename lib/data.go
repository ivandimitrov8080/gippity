package lib

import (
	"github.com/sashabaranov/go-openai"
)

func SplitInfoIntoSubsections(info string, maxTokens int) []string {
	ans := []string{}
	ans = append(ans, breakTextEvenly(info, maxTokens)...)
	return ans
}

func tokenCount(text string) int {
	return CountTokens(text, openai.GPT3Dot5Turbo0301)
}

func breakTextEvenly(text string, maxTokens int) []string {
	if tokenCount(text) < maxTokens {
		return []string{text}
	}
	split := splitInTwo(text)
	ans := []string{split[0], split[1]}
	var temp []string
	for {
		temp = []string{}
		for _, v := range ans {
			split := splitInTwo(v)
			temp = append(temp, split[:]...)
		}
		ans = temp
		if tokenCount(ans[0]) < maxTokens {
			break
		}
	}
	return ans
}

func splitInTwo(text string) [2]string {
	return [2]string{
		text[0 : len(text)/2],
		text[len(text)/2:],
	}
}
