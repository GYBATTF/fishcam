all: fishcam

fishcam:
	go mod download
	go build -o fishcam ./src/.
