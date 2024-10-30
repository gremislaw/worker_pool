all: build

build:
	go mod download
	go build -o ./bin/worker_pool

rebuild:
	make build

run: build
	./bin/worker_pool

test:
	go test ./...

style:
	go fmt ./...

clean:
	rm -rf ./bin out.txt
