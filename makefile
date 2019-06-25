.DEFAULT_GOAL := default

.PHONY: docker binary

binary:
	CGO_ENABLED=0 go build

docker: Dockerfile
	docker build . -t discord-url-shortener

default: binary
