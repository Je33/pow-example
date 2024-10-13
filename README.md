# Example of using PoW concept to protect server from DDoS attacks

## Description

This example shows how to use PoW concept to protect server from DDoS attacks. The algorithm is based on [SHA-3 family](https://en.wikipedia.org/wiki/SHA-3) and [PoW hashcash](https://en.wikipedia.org/wiki/Hashcash) because it is the lightweight, simplest and modern proof-of-work algorithm, based firstly on CPU intensive proof of work, but solution can be extended to other PoW algorithms, based on requirements.

## Prerequisites

- Docker
- Golang

## Usage

To start server and client with docker, you can use the following commands:
```bash
docker-compose up --build
```
Or it can be started just in separate processes on one host:
```bash
go build -o ./build/server ./cmd/server/main.go
go build -o ./build/client ./cmd/client/main.go
```
And start in different terminals:
```bash
./build/server
```
```bash
./build/client
```

## Tests

```bash
go test ./...
```

## Makefile

- `build` - build server and client
```bash
make build
```
- `run` - start server and client in docker
```bash
make dc
```
- `test` - run tests server and client with coverage report
```bash
make test
```
- `lint` - run linter
```bash
make lint
```


## References

- [PoW concept](https://en.wikipedia.org/wiki/Proof_of_work)
- [PoW hashcash](https://en.wikipedia.org/wiki/Hashcash)
- [SHA-3 family](https://en.wikipedia.org/wiki/SHA-3)
- [DDoS attack](https://en.wikipedia.org/wiki/DDoS_attack)