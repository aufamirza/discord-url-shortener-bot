FROM golang:1.12.6 as build
WORKDIR /build
COPY / .
ENV CGO_ENABLED=0
RUN go build -o discord-url-shortener-bot
FROM alpine:3.10.0
COPY --from=build /build/discord-url-shortener-bot .
ENTRYPOINT ./discord-url-shortener-bot
