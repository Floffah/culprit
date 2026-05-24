default:
	@just --list

build:
    go build -trimpath -ldflags="-w -s" -o dist/culprit cmds/maculprit/maculprit.go

test:
	go test -v -cover ./...

lint:
	golangci-lint run