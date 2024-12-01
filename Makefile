build:
	GOARCH=arm64 GOOS=linux go build -o main main.go

zip:
	zip deployment.zip main bootstrap ffmpeg