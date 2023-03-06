.DEFAULT_GOAL := build

app.name := dumpcft
app.repo := github.com/scottbrown/$(app.name)

.PHONY: build
build:
	go build -o .build/$(app.name) $(app.repo)/cmd

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
