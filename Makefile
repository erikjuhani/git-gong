GOBUILD=go build
GOTEST=go test
GOCOVER=go tool cover
ARM64_ARCH=GOARCH=arm64
AMD64_ARCH=GOARCH=amd64

build:
	@mkdir -p bin
	@GOOS=darwin $(ARM64_ARCH) $(GOBUILD) -tags static,system_libgit2 -o bin/gong-darwin-x86_64 main.go
	@# GOOS=linux $(AMD64_ARCH) $(GOBUILD) -tags static,system_libgit2 -o bin/gong-linux-x86_64 main.go
	@# GOOS=windows $(AMD64_ARCH) $(GOBUILD) -tags static,system_libgit2 -o bin/gong.exe main.go

test:
	@$(GOTEST) -tags static,system_libgit2 ./...

coverage:
	@$(GOTEST) -tags static,system_libgit2 -coverprofile=coverage.out ./...
	@$(GOCOVER) -html=coverage.out
	@rm coverage.out

.PHONY: build test coverage
