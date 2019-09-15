# Alpago - An Alpaca Bot
Part of [DumbergerL's Bot Challenge](https://github.com/DumbergerL/alpaca-server).

It is written in pure Go and only uses Go's standard library.

# Usage

## Prerequisites
- Install and setup [Go](https://golang.org/dl/) if you haven't yet. That's all folks.

## Playing
Either run
`go build main.go`
and then execute resulting binary or just run
`go run main.go`.

By default alpago will connect to `http://localhost:3000`, set `localhost` for callbacks, spawn four players and use 1 as drop limit, this can be changed with the `-host`, `-port`, `-ip`, `-players` and `-limit` flags.

E.g. to connect to `http://10.0.0.1:1234` with one player run `go run main.go -host=10.0.0.1 -port=1234 -ip=<your ip address> -players=1 -limit=2`.
