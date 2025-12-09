# KV-Store

A minimal, key-value cache written in Go.
Runs a lightweight local daemon and exposes a smple command-based protocol over a UNIX socket.

Designed for personal automation, scripting, and integration with other software.

## Features

- In-memory index
- `GET`/`PUT`/`DELETE` operations
- Tombstone-based deletion
- Simple CLI client
- Fully local, zero external dependencies

## Requirements
- Go 1.20+
- macOS / Linux (Linux untested; uses UNIX sockets)

## Installation

Clone this repo.

```bash
git clone https://github.com/arianfarid/kvstore
cd kvstore```

By default, the kv file is located in the repo.

There are two binaries to build. I chose to build them local to the project, though you may build in your desired location.
Build the daemon:
```bash
go build -o ./bin/kvd ./cmd/main.go
```

Build the CLI tool:
```bash
go build -o ./bin/kv ./cmd/kv/main.go
```

## Usage

First, run the daemon:
```bash
./bin/kvd 
```
This creates and listens on the UNIX socket `/tmp/kvstore.sock`.

Store a value:
```bash
./bin/kv PUT test This is a test value #outputs: OK
```

Retrieve a value: 
```bash
./bin/kv GET test #outputs: VALUE this is a test value
```
Delete a value:
```bash
./bin/kv DELETE test #outputs: OK
```


## Features Planned

- Log compaction
- Offset-based value lookup
- Optional encryption
