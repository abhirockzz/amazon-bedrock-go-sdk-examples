package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/abhirockzz/amazon-bedrock-go-inference-params/claude"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

const defaultRegion = "us-east-1"

const (
	claudePromptFormat = "\n\nHuman:%s\n\nAssistant:"
	claudeV2ModelID    = "anthropic.claude-v2" //https://docs.aws.amazon.com/bedrock/latest/userguide/model-ids-arns.html
)

const prompt = `<paragraph> 
"In 1758, the Swedish botanist and zoologist Carl Linnaeus published in his Systema Naturae, the two-word naming of species (binomial nomenclature). Canis is the Latin word meaning "dog", and under this genus, he listed the domestic dog, the wolf, and the golden jackal."
</paragraph>

Please rewrite the above paragraph to make it understandable to a 5th grader.

Please output your rewrite in <rewrite></rewrite> tags.`

func main() {

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = defaultRegion
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	brc := bedrockruntime.NewFromConfig(cfg)

	payload := claude.Request{
		Prompt:            fmt.Sprintf(claudePromptFormat, prompt),
		MaxTokensToSample: 2048,
		Temperature:       0.5,
		TopK:              250,
		TopP:              1,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	output, err := brc.InvokeModelWithResponseStream(context.Background(), &bedrockruntime.InvokeModelWithResponseStreamInput{
		Body:        payloadBytes,
		ModelId:     aws.String(claudeV2ModelID),
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		log.Fatal("failed to invoke model: ", err)
	}

	resp, err := processStreamingOutput(output, func(ctx context.Context, part []byte) error {
		fmt.Print(string(part))
		return nil
	})

	if err != nil {
		log.Fatal("streaming output processing error: ", err)
	}

	fmt.Println("\n====== response from LLM ======\n", resp.Completion)

}

type StreamingOutputHandler func(ctx context.Context, part []byte) error

func processStreamingOutput(output *bedrockruntime.InvokeModelWithResponseStreamOutput, handler StreamingOutputHandler) (claude.Response, error) {

	var combinedResult string

	for event := range output.GetStream().Events() {
		switch v := event.(type) {
		case *types.ResponseStreamMemberChunk:

			//fmt.Println("payload", string(v.Value.Bytes))

			var resp claude.Response
			err := json.NewDecoder(bytes.NewReader(v.Value.Bytes)).Decode(&resp)
			if err != nil {
				fmt.Println("Error unmarshalling JSON:", err)
				continue
			}

			handler(context.Background(), []byte(resp.Completion))
			combinedResult += resp.Completion

		case *types.UnknownUnionMember:
			fmt.Println("unknown tag:", v.Tag)

		default:
			fmt.Println("union is nil or unknown type")
		}
	}

	//fmt.Println("final result", combinedResult)

	resp := claude.Response{}
	resp.Completion = combinedResult

	return resp, nil
}
