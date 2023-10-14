# iho

a dumb TCP tunnel

### Usage

```sh
iho [flags] [mode]
```

There are two modes, `server` & `client`.

```sh
go run ./cmd/iho -listen :3333 server
```

```sh
go run ./cmd/iho -to :8000 -remote :3333 client
```

### Features

[x] no authentication
