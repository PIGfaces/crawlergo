.PHONY: wire
wire:
	cd pkg && wire

.PHONY: build
# build
build:
	rm -rf bin && mkdir bin bin/linux-amd64 bin/linux-arm64 bin/darwin-amd64 bin/darwin-arm64 \
	&& CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/darwin-arm64/ ./... \
	&& CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/darwin-amd64/ ./... \
	&& CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/linux-arm64/ ./... \
	&& CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/linux-amd64/ ./...