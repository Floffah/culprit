default:
	@just --list

build:
    go build -trimpath -ldflags="-w -s" -o dist/culprit cmds/culprit/culprit.go

test:
	go test -v -cover ./...

lint:
	golangci-lint run

format:
	gofmt -s -w .
	golangci-lint fmt