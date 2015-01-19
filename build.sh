#!/bin/bash
export GOPATH=`pwd`
go build main && mv main ./bin
go build client && mv client ./bin
go build testclient && mv testclient ./bin
go build testserver && mv testserver ./bin
