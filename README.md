# go-grpc-downgrade
This package donwgrades generated golang grpc package version.

## Motivation
[protoc-gen-go-grpc](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc) from version 6 to 7 changed file generated adding xyz_grpc.pb.go.

This package was created to mantain backward compatibility.

## Features
* Downgrade Golang grpc generated package version 7 to 6
```go
const _ = grpc.SupportPackageIsVersion7
```
to
```go
const _ = grpc.SupportPackageIsVersion6
```

## Install

```shell
go get github.com/fakiot/go-grpc-downgrade
```


## Quick start
Generate Golang grpc code using [protoc-gen-go-grpc](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc) then add the following go generate directive
```go
//go:generate go-grpc-downgrade -pb xyz.pb.go -grpc xyz_grpc.pb.go -o xyz.pb.go
```