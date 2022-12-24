package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

type ImageRequest struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type ImageResponse struct {
	Created int `json:"created"`
	Data    []struct {
		URL string `json:"url"`
	} `json:"data"`
}

var OPENAI_API_KEY string

func main() {
	OPENAI_API_KEY := os.Getenv("API_KEY")

	if OPENAI_API_KEY == "" {
		fmt.Println("no api key found")
		panic("env not correct")
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter an imge description : ")
	prompt, _ := reader.ReadString('\n')
	GenerateImage(prompt)
}

func GenerateImage(prompt string) {

	requestBody, err := json.Marshal(ImageRequest{
		Prompt: prompt,
		Size:   "256x256",
		N:      1,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", bytes.NewBuffer(requestBody))
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

	var image ImageResponse
	if err := json.NewDecoder(resp.Body).Decode(&image); err != nil {
		fmt.Println(err)
		return
	}

	parsedURL, err := url.Parse(image.Data[0].URL)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = exec.Command("open", parsedURL.String()).Start()
	if err != nil {
		fmt.Println(err)
		return
	}
}
