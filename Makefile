build:
	GOARCH=arm64 GOOS=linux go build -o main function-video2gif/main.go

zip:
	zip deployment.zip main bootstrap bin/ffmpeg

quick-deploy:
	aws lambda update-function-code --function-name GoFunction --zip-file fileb://deployment.zip