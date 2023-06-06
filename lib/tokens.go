package lib

import (
	"fmt"

	"github.com/pkoukk/tiktoken-go"
)

func CountTokens(message string, model string) (num_tokens int) {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("EncodingForModel: %v", err)
		fmt.Println(err)
		return
	}
	return len(tkm.Encode(message, nil, nil))
}
