build-ffmpeg-version:
	GOARCH=arm64 GOOS=linux go build -o deployment/ffmpeg-version/main function-ffmpeg-version/main.go
	cd deployment/ffmpeg-version && zip deployment.zip main ../../bootstrap ../../bin/ffmpeg

build-video2gif:
	GOARCH=arm64 GOOS=linux go build -o deployment/video2gif/main function-video2gif/main.go
	cd deployment/video2gif && zip deployment.zip main ../../bootstrap ../../bin/ffmpeg

build-all: build-ffmpeg-version build-video2gif

quick-deploy-ffmpeg-version:
	aws lambda update-function-code --function-name FfmpegVersionFunction --zip-file fileb://deployment/ffmpeg-version/deployment.zip

quick-deploy-video2gif:
	aws lambda update-function-code --function-name Video2GifFunction --zip-file fileb://deployment/video2gif/deployment.zip
