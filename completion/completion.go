package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type CompletionRequest struct {
	Prompt      string  `json:"prompt"`
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	N           int     `json:"n"`
	Temperature float32 `json:"temperature"`
}

type CompletionResponse struct {
	Id      string `json:"id"`
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	Choices []struct {
		Text  string `json:"text"`
		Index int    `json:"index"`
	} `json:"choices"`
}

const (
	OPENAI_API_KEY = ""
)

func main() {
	OPENAI_API_KEY := os.Getenv("API_KEY")
	// var prompt string
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a line to predict: ")
	prompt, _ := reader.ReadString('\n')

	model := "text-davinci-002"
	requestBody, err := json.Marshal(CompletionRequest{
		Prompt:      prompt,
		Model:       model,
		MaxTokens:   1024,
		N:           4,
		Temperature: 0.5,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+OPENAI_API_KEY)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	var response CompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println(err)
		return
	}
	// message := response.Choices[0].Text
	// promptPrefix := strings.TrimRight(prompt, "\n")
	// messageTrimmed := strings.TrimPrefix(message, "\n")
	// fmt.Println(promptPrefix + " " + messageTrimmed + "\n")

	for i := 0; i < len(response.Choices); i++ {
		fmt.Println(response.Choices[i].Index)
		fmt.Println(response.Choices[i].Text)
	}

}
