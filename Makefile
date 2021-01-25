GOBUILD=go build
GOTEST=go test
GOCOVER=go tool cover
AMD64_ARCH=GOARCH=amd64

build:
	mkdir -p bin
	GOOS=darwin $(AMD64_ARCH) $(GOBUILD) -installsuffix cgo -tags static -o bin/gong-darwin-x86_64 main.go
	# GOOS=linux $(AMD64_ARCH) $(GOBUILD) -installsuffix cgo -tags static -o bin/gong-linux-x86_64 main.go
	# GOOS=windows $(GOBUILD) -installsuffix cgo -tags static -o bin/gong.exe main.go

test:
	$(GOTEST) -tags static ./...

coverage:
	@$(GOTEST) -tags static -coverprofile=coverage.out ./...
	@$(GOCOVER) -html=coverage.out
	@rm coverage.out

.PHONY: build test coverage
