.PHONY: build clean update fmt test test-small test-medium coverage

build:
	go build -v .

clean:
	go clean
	go mod tidy

update:
	go get -u

fmt:
	go fmt ./...

lint:
	staticcheck ./...

test:
	go test -v -tags=small,medium ./...

test-small:
	go test -v -tags=small ./...

test-medium:
	go test -v -tags=medium ./...

coverage:
	go test ./... -tags=small,medium -covermode=count -coverprofile=c.out
	go tool cover -html=c.out -o coverage.html
