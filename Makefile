all: run

run:
	go build
	./fake-server -port=8080
	git clean -fd
