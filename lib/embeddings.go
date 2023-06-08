package lib

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/neurosnap/sentences/english"
	openai "github.com/sashabaranov/go-openai"
)

const batchSize = 1500

func CreateEmbeddings(client *openai.Client, text string, maxTokens int) map[string][]float32 {
	fmt.Println("Started embedding")
	tokenizer, err := english.NewSentenceTokenizer(nil)
	if err != nil {
		fmt.Println(err)
	}
	sentences := tokenizer.Tokenize(text)
	embeds := []openai.Embedding{}
	sentenceCount := len(sentences)
	for i := 0; i < sentenceCount; i += batchSize {
		end := i + batchSize
		if end > sentenceCount {
			end = sentenceCount
		}
		var batch []string
		for i := 0; i < sentenceCount; i++ {
			batch = append(batch, sentences[i].Text)
		}
		fmt.Println("Embedding batch: ", i, " :: ", end)
		resp, err := client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
			Input: batch,
			Model: openai.AdaEmbeddingV2,
		})

		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		vectors := resp.Data
		embeds = append(embeds, vectors...)
	}
	if len(embeds) == sentenceCount {
		m := map[string][]float32{}
		for i, v := range sentences {
			m[v.Text] = embeds[i].Embedding
		}
		return m
	}
	return map[string][]float32{}
}

func CosineDistance(v1 []float32, v2 []float32) float32 {
	var operand float32 = 0
	var d1 float32 = 0
	var d2 float32 = 0
	for i, v := range v1 {
		operand += v * v2[i]
		d1 += v * v
		d2 += v2[i] * v2[i]
	}
	return operand / (float32(math.Sqrt(float64(d1))) * float32(math.Sqrt(float64(d2))))
}

func relatedness(v1 []float32, v2 []float32) float32 {
	return 1 - CosineDistance(v1, v2)
}

func TopMatches(questionEmbeddings map[string][]float32, knowledgeEmbeddings map[string][]float32, maxTokens int) []string {
	ans := []string{}
	m := map[string]float32{}
	var qe []float32
	type kv struct {
		Key   string
		Value float32
	}
	for _, v := range questionEmbeddings {
		qe = v
		break
	}
	for k, v := range knowledgeEmbeddings {
		m[k] = relatedness(qe, v)
	}
	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	l := len(ss)
	for i := 0; maxTokens > CountTokens(strings.Join(ans, ""), openai.GPT3Dot5Turbo0301) && i < l; i++ {
		ans = append(ans, ss[i].Key)
	}
	return ans
}
