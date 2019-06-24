# DISCORD-URL-SHORTENER-BOT

This bot will automatically generate short redirect URL's and replace the original message in any Discord channel it has permissions to read.

Format `<PROTOCOL>://<HOSTNAME>/<LONG_PATH_AND_QUERY_PARAMS>` => `<PROTOCOL>://<HOSTNAME>/<SHORT_PATH>`.

E.g. `https://github.com/fraserdarwent/discord-url-shortener-bot` => `https://<DOMAIN>/NHhG`.

### Persistence

Persistence will be a local file backend in the form of a `data.json` in the app directory.
More persistence backends may be added.
