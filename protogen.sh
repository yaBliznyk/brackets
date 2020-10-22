#!/bin/bash

protoc -I . \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--go_out ./ --go_opt paths=source_relative \
--go-grpc_out ./ --go-grpc_opt paths=source_relative \
proto/brackets.proto

protoc -I . --grpc-gateway_out ./ \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--grpc-gateway_opt logtostderr=true \
--grpc-gateway_opt paths=source_relative \
proto/brackets.proto