.PHONY: run build test lint docker clean

run:
	cd cmd/api-service && ENV=LOCAL PORT=3001 VERSION=VERSION go run .

build:
	go build -o bin/api-service ./cmd/api-service

test:
	go test ./... -v -cover

lint:
	golangci-lint run

docker:
	docker build -t go-rest-api-template .

clean:
	rm -rf bin/
