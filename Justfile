set quiet

alias b := build

build:
    go build ./... && go vet ./...