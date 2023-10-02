package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/abhirockzz/amazon-bedrock-go-inference-params/stabilityai"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	stableDiffusionXLModelID = "stability.stable-diffusion-xl-v0" //https://docs.aws.amazon.com/bedrock/latest/userguide/model-ids-arns.html
)

type Request struct {
	Input string `json:"input"`
}

type Response struct {
	Output string `json:"output"`
}

func handler(req Request) (Response, error) {
	prompt := req.Input

	log.Println("input", prompt)

	payload := stabilityai.Request{
		TextPrompts: []stabilityai.TextPrompt{{Text: prompt}},
		CfgScale:    10,
		Seed:        0,
		Steps:       50,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	output, err := brc.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		Body:        payloadBytes,
		ModelId:     aws.String(stableDiffusionXLModelID),
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		log.Fatal("failed to invoke model: ", err)
	}

	var resp stabilityai.Response

	err = json.Unmarshal(output.Body, &resp)

	if err != nil {
		log.Fatal("failed to unmarshal", err)
	}

	decoded, err := resp.Artifacts[0].DecodeImage()

	if err != nil {
		log.Fatal("failed to decode base64 response", err)
	}

	outputFile := fmt.Sprintf("output-%d.jpg", time.Now().UnixMilli())
	s3File := fmt.Sprintf("s3://%s/%s", targetBucket, outputFile)

	fmt.Println("image generated. uploading to", s3File)

	uploadReq := &s3.PutObjectInput{
		Bucket: aws.String(targetBucket),
		Key:    aws.String(outputFile),
		Body:   bytes.NewReader(decoded),
	}
	_, err = s3Client.PutObject(context.Background(), uploadReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("successfully uploaded image to s3 bucket", targetBucket)

	return Response{Output: fmt.Sprintf("s3://%s/%s", targetBucket, outputFile)}, nil
}

const defaultRegion = "us-east-1"

var brc *bedrockruntime.Client
var targetBucket string
var s3Client *s3.Client

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

	targetBucket = os.Getenv("BUCKET_NAME")
	if targetBucket == "" {
		log.Fatal("missing environment variable BUCKET_NAME")
	}

	fmt.Println("target S3 bucket", targetBucket)

	s3Client = s3.NewFromConfig(cfg)

}

func main() {
	lambda.Start(handler)
}
