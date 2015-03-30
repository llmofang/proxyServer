#!/bin/bash
export GOPATH=`pwd`
if [ ! -d  `pwd`"/src/github.com/hoisie/redis" ]
then	
	echo '==> Getting dependencies (hoisie/redis) ...'
	go get  github.com/hoisie/redis
fi
if [ ! -d  `pwd`"/src/github.com/op/go-logging" ]
then	
	echo '==> Getting dependencies (op/go-logging) ...'
	go get  github.com/op/go-logging
fi
echo '==> building'
go build main && mv main ./bin