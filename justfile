alias gen := generate

default:
	@just --list

build:
    go build -trimpath -ldflags="-w -s" -o dist/culprit cmds/culprit/culprit.go

test:
	go test -v -cover ./...

test-integration:
	go run cmds/culprit/culprit.go --verbose clean soft --force -o cleanup-soft.sh
	go run cmds/culprit/culprit.go --verbose clean hard --force -o cleanup-hard.sh
	go run cmds/culprit/culprit.go --verbose clean culprit --force -o cleanup-culprit.sh


lint:
	golangci-lint run

format:
	gofmt -s -w .
	golangci-lint fmt

generate:
	go generate ./...