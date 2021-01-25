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
	$(GOTEST) ./...

coverage:
	@$(GOTEST) -coverprofile=coverage.out ./...
	@$(GOCOVER) -html=coverage.out
	@rm coverage.out

.PHONY: test coverage
