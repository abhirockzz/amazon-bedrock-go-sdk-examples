package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/abhirockzz/amazon-bedrock-go-inference-params/claude"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type Request struct {
	Input string `json:"input"`
}

type Response struct {
	Output string `json:"output"`
}

const claudePromptFormat = "\n\nHuman:%s\n\nAssistant:"

func handler(req Request) (Response, error) {
	prompt := req.Input

	log.Println("input", prompt)

	payload := claude.Request{
		Prompt:            fmt.Sprintf(claudePromptFormat, prompt),
		MaxTokensToSample: 100,
		//Temperature:       0.7,
		//TopP:              0.95,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	output, err := brc.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		Body:        payloadBytes,
		ModelId:     aws.String("anthropic.claude-v2"),
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		log.Fatal("failed to invoke model", err)
	}

	log.Println("raw response ", string(output.Body))

	var resp claude.Response

	err = json.Unmarshal(output.Body, &resp)

	if err != nil {
		log.Fatal("failed to unmarshal", err)
	}

	fmt.Println("response from LLM\n", resp.Completion)

	return Response{Output: resp.Completion}, nil
}

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
	lambda.Start(handler)
}
