#!/bin/bash
export GOPATH=`pwd`
go build main && mv main ./bin
go build client && mv client ./bin
