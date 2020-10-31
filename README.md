# descrob

A less click intensive way to unscrobble tracks from [Last.FM](https://last.fm). Provides a package to programmatically unscrobble and a TUI app to navigate scrobble history an unscrobble directly from the shell.

## Getting started

Still in early stages of development, so only distributed via go toolchain for now.

Requires a Last.FM developer account and an application API key (see [here](https://www.last.fm/api)).

```shell
go get github.com/Morley93/descrob/...

LASTFM_USERNAME=<username> LASTFM_PASSWORD=<password> LASTFM_API_KEY=<api_key> $GOPATH/bin/descrob
```

## Motivation

After experiencing some double-scrobbles and getting frustrated while dutifully removing them via the web UI, I wanted a way to make selective unscrobbling easier.

The Last.FM API doesn't support unscrobbling tracks, so this tool achieves it by emulating a web session.