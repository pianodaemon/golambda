GOCMD=go
GOBUILD=$(GOCMD) build -ldflags="-w -s"

LAMBDA_EXE=main

all: build

build: fmt
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	$(GOBUILD) -o $(LAMBDA_EXE) lambda/main.go

fmt:
	$(GOCMD) fmt ./...

clean:
	rm -f $(LAMBDA_EXE)
