package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"strings"
	"time"
)

type Response struct {
	Message string `json:"message"`
}

var (
	s3Client *s3.Client
)

func init() {
	// Initialize the S3 client outside the handler, during the init phase
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
	if s3Client == nil {
		log.Fatalf("Failed to initialize S3 client")
	}
}

func getFile(filePath string, bucketName string) (*s3.GetObjectOutput, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filePath),
	}

	res, err := s3Client.GetObject(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	if res == nil {
		return nil, fmt.Errorf("received nil response from GetObject")
	}

	return res, nil
}

func uploadFile(filePath string, bucketName string) (*s3.PutObjectOutput, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filePath),
		Body:   strings.NewReader(time.Now().String()),
	}

	res, err := s3Client.PutObject(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	if res == nil {
		return nil, fmt.Errorf("received nil response from PutObject")
	}

	fmt.Printf("%+v\n", res.ResultMetadata)
	return res, nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	uploaded, err := uploadFile("something", "iamtestingmys3")
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	fmt.Printf("%+v\n", uploaded)

	gotFile, err := getFile("something", "iamtestingmys3")
	if err != nil {
		log.Printf("Error getting file: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	fmt.Printf("%+v\n", gotFile)

	response := Response{
		Message: fmt.Sprintf("Successfully uploaded, metadata: %v", gotFile.Body),
	}

	body, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}
