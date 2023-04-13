GOCMD=go
GOBUILD=$(GOCMD) build -ldflags="-w -s"

LAMBDA		= main
LAMBDA_EXE	= $(LAMBDA)
LAMBDA_GO	= $(LAMBDA).go
LAMBDA_ZIP	= $(LAMBDA).zip

all: format zip

build:
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	$(GOBUILD) -o $(LAMBDA_EXE) lambda/$(LAMBDA_GO)

zip: build
	zip $(LAMBDA_ZIP) $(LAMBDA_EXE)

format:
	$(GOCMD) fmt ./...

clean:
	rm -f $(LAMBDA_EXE) $(LAMBDA_ZIP)
