# go-webterm: Terminal over the Browser

go-webterm is a [xterm.js](https://xtermjs.org/) Go backend for learning (and fun!) purposes

## Requirements

- Go >= 1.20

## Build and Usage

Running `make build` will generate the binary `bin/webterm`

```bash
bin/webterm
2023/03/05 19:35:09 INFO Listening on localhost:8000
2023/03/05 19:35:43 INFO Received connection from: 127.0.0.1:53342
2023/03/05 19:35:43 INFO Starting TTY
2023/03/05 19:35:43 INFO Waiting
2023/03/05 19:35:43 INFO Sending bytes from PTY to Websocket
2023/03/05 19:35:43 INFO Sending bytes from PTY to Websocket
2023/03/05 19:35:45 INFO Copying bytes from Websocket to TTY
2023/03/05 19:35:45 INFO Sending bytes from PTY to Websocket
```

## Tests

```bash
make test
?       github.com/pabloxio/go-webterm/cmd/webterm      [no test files]
?       github.com/pabloxio/go-webterm/webterm  [no test files]
ok      github.com/pabloxio/go-webterm/handlers 0.272s  coverage: 12.5% of statements
```
