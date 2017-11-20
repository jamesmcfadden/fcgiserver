# fcgiserver

:fire: **Do not use in a production environment** :fire:

## Installation

`go get github.com/jamesmcfadden/fcgiserver`

## Pre requisites

1. Ensure `$GOPATH/bin` is in your path
2. Ensure you have a FCGI implementation running on `127.0.0.1:9000`
3. Run the following command:

`fcgiserver -l localhost:8000 -r /workspace/go-server`

Flags:

`-l` Listen on

`-r` Document root
