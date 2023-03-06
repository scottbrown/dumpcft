.DEFAULT_GOAL := build

app.name := dumpcft
app.repo := github.com/scottbrown/$(app.name)

build:
	go build -o .build/$(app.name) $(app.repo)/cmd

test:
	go test ./...

fmt:
	go fmt ./...
