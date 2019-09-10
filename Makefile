PORT ?= 8080

build:
	go build -o bin/lumberman-web-client .

build-linux64:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/lumberman-web-client .

run:
	./bin/lumberman-web-client

dev: build
	./bin/lumberman-web-client -port ${PORT}

.PHONY: build build-linux64 run dev
