run:
	go build
	./fake-server
	git clean -fd
