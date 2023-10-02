package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/abhirockzz/amazon-bedrock-go-inference-params/claude"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

const defaultRegion = "us-east-1"

var brc *bedrockruntime.Client

func init() {

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = defaultRegion
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	brc = bedrockruntime.NewFromConfig(cfg)

}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// initial conversation state
	chatInput := Chat{}

	for {
		fmt.Print("\nEnter your message: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		chatInput.CurrentMessage = input

		response, err := chatInput.send()

		if err != nil {
			log.Fatal(err)
		}

		chatInput.History = chatInput.History + response

		fmt.Println("\n--- Response ---")
		fmt.Println(response)
	}
}

type Chat struct {
	History        string
	CurrentMessage string
}

func (ci Chat) getPayload() string {
	return ci.History + fmt.Sprintf(claudePromptFormat, ci.CurrentMessage)
}

const claudePromptFormat = "\n\nHuman:%s\n\nAssistant:"

func (ci Chat) send() (string, error) {

	msg := ci.getPayload()

	//log.Println("sending message", msg)

	payload := claude.Request{Prompt: msg, MaxTokensToSample: 2048}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	output, err := brc.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		Body:        payloadBytes,
		ModelId:     aws.String("anthropic.claude-v2"),
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		return "", err
	}

	var resp claude.Response

	err = json.Unmarshal(output.Body, &resp)

	if err != nil {
		return "", err
	}

	return resp.Completion, nil
}
