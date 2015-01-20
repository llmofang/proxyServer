#!/bin/bash
export GOPATH=`pwd`
go build main && mv main ./bin
