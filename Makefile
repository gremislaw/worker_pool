all: build

build:
	go mod download
	go build -o ./bin/worker_pool

rebuild:
	make build

run: build
	./bin/worker_pool

test:
	go test -v ./pool
	make clean

style:
	go fmt ./...

clean:
	rm -rf ./bin *.txt pool/*.txt worker_pool
