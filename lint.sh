#!/bin/sh
go mod tidy
golangci-lint run --fix