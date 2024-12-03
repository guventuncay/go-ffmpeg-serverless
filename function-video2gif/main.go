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
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"log"
	"os"
)

type Response struct {
	Message string `json:"message"`
}

var (
	s3Client *s3.Client
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
}

func getFile(filePath string, bucketName string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filePath),
	}

	res, err := s3Client.GetObject(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file body: %w", err)
	}

	return data, nil
}

func convertToGif(video []byte) (string, error) {
	inputFile := "/tmp/input.mp4"
	outputFile := "/tmp/output.gif"

	err := os.WriteFile(inputFile, video, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write video to file: %w", err)
	}

	err = ffmpeg.Input(inputFile, ffmpeg.KwArgs{"ss": "1"}).
		Output(outputFile, ffmpeg.KwArgs{"s": "320x240", "pix_fmt": "rgb24", "t": "3", "r": "3"}).
		OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		return "", fmt.Errorf("failed to convert video to GIF: %w", err)
	}

	return outputFile, nil
}

func uploadFile(filePath string, bucketName string) (*s3.PutObjectOutput, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for upload: %w", err)
	}
	defer file.Close()

	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("output.gif"),
		Body:   file,
	}

	res, err := s3Client.PutObject(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return res, nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	videoData, err := getFile("video.mp4", "iamtestingmys3")
	if err != nil {
		log.Printf("Error getting video: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error downloading video"}, err
	}

	gifPath, err := convertToGif(videoData)
	if err != nil {
		log.Printf("Error converting video: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error converting video to GIF"}, err
	}

	_, err = uploadFile(gifPath, "iamtestingmys3")
	if err != nil {
		log.Printf("Error uploading GIF: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error uploading GIF"}, err
	}

	response := Response{Message: "GIF created and uploaded successfully"}
	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(body)}, nil
}

func main() {
	lambda.Start(handler)
}
