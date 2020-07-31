#!/usr/bin/env sh

protoc --proto_path=proto --proto_path=third_party --go_out=plugins=grpc:proto file.proto