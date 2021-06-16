BLD_NAME:=ptt-crawler

.PHONY: build 
build: cmd/$(BLD_NAME)/main.go 
	go build -o bin/$(BLD_NAME) $<

.PHONY: build-image
build-image: Dockerfile
	podman build -f Dockerfile -t ptt-crawler:latest
