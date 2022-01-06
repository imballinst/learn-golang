# tutorial-greetings

This learning section is based on https://go.dev/doc/tutorial/create-module.

## Requirements

This requires Go version at least `1.17.3`.

## How to Run

```shell
go run ./hello.go
# map[Darrin:Hi, Darrin. Welcome! Gladys:Hail, Gladys! Well met! Samantha:Hi, Samantha. Welcome!]
```

## How to Test

```shell
cd greetings
go test
# PASS
# ok  	example.com/greetings	0.001
```

## Building

```shell
./build.sh
./bin/hello
# map[Darrin:Hi, Darrin. Welcome! Gladys:Hail, Gladys! Well met! Samantha:Hi, Samantha. Welcome!]
```
