.DEFAULT_GOAL := build

app.name := dumpcft
app.repo := github.com/scottbrown/$(app.name)

pwd := $(shell pwd)

dist.dir := $(pwd)/.dist
build.dir := $(pwd)/.build

.PHONY: build
build:
	go build -o $(build.dir)/$(app.name) $(app.repo)/cmd

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: check
check: sast vet vuln

.PHONY: sast
sast:
	gosec ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: vuln
vuln:
	govulncheck ./...

.PHONY: clean
clean:
	rm -rf $(build.dir) $(dist.dir)

.PHONY: release
release:
	GOOS=linux GOARCH=amd64 go build -o $(build.dir)/linux-amd64/dumpcft $(app.repo)/cmd
	GOOS=linux GOARCH=arm64 go build -o $(build.dir)/linux-arm64/dumpcft $(app.repo)/cmd
	GOOS=darwin GOARCH=amd64 go build -o $(build.dir)/darwin-amd64/dumpcft $(app.repo)/cmd
	GOOS=darwin GOARCH=arm64 go build -o $(build.dir)/darwin-arm64/dumpcft $(app.repo)/cmd
	GOOS=windows GOARCH=amd64 go build -o $(build.dir)/windows-amd64/dumpcft $(app.repo)/cmd
	GOOS=windows GOARCH=arm64 go build -o $(build.dir)/windows-arm64/dumpcft $(app.repo)/cmd
	mkdir -p $(dist.dir)
	tar cfz $(dist.dir)/dumpcft.linux-amd64.tar.gz -C $(build.dir)/linux-amd64 .
	tar cfz $(dist.dir)/dumpcft.linux-arm64.tar.gz -C $(build.dir)/linux-arm64 .
	tar cfz $(dist.dir)/dumpcft.darwin-amd64.tar.gz -C $(build.dir)/darwin-amd64 .
	tar cfz $(dist.dir)/dumpcft.darwin-arm64.tar.gz -C $(build.dir)/darwin-arm64 .
	tar cfz $(dist.dir)/dumpcft.windows-amd64.tar.gz -C $(build.dir)/windows-amd64 .
	tar cfz $(dist.dir)/dumpcft.windows-arm64.tar.gz -C $(build.dir)/windows-arm64 .
