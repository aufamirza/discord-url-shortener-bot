.DEFAULT_GOAL := default

.PHONY: docker binary

#build the binary
binary:
	CGO_ENABLED=0 go build

#build the binary into a Docker image
docker: Dockerfile
	docker build . -t discord-url-shortener

default: binary
