# DISCORD-URL-SHORTENER-BOT
This bot will automatically generate short redirect URL's and replace the original message in any Discord channel it has permissions to read.

Format `<PROTOCOL>://<HOSTNAME>/<LONG_PATH_AND_QUERY_PARAMS>` => `<PROTOCOL>://<HOSTNAME>/<SHORT_PATH>`.

E.g. `https://github.com/fraserdarwent/discord-url-shortener-bot` => `https://<DOMAIN>/NHhG`.

## Persistence
Persistence will be a local file backend in the form of a `data.json` in the app directory.
More persistence backends may be added.

## Build The Bot To Run Your Own Instance

### Binary
```shell script
make binary
```

### Docker
```shell script
make docker
```

## Run Your Own Instance
Build the binary
```shell script
DISCORD_BOT_TOKEN=<YOUR_TOKEN_HERE> \
DISCORD_BOT_PROTOCOL=<HTTP || HTTPS> \
DISCORD_BOT_HOST=<YOUR_HOST_HERE> \
DISORD_BOT_PORT=<WEB_SERVER_PORT> \
go run main.go
```
## Add My Instance Of The Bot To Your Server
https://discordapp.com/api/oauth2/authorize?client_id=592857117548085288&permissions=10240&redirect_uri=https%3A%2F%2Fdiscordapp.com%2Foauth2%2Fauthorize&scope=bot
