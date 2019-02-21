GOARCH=amd64
GOOS=linux

build:
	GOOS=${GOOS} GOARCH=${GOARCH} go build .;
run:
	./tracking_server;
clean:
	go clean;
