#!/bin/bash
export GOPATH=`pwd`
if [ -e  `pwd`+"/src/github.com/astaxie/redis" ]
then	
	echo 'getting github.com/hoisie/redis'
	go get  github.com/hoisie/redis
fi
if [ -e  `pwd`+"/src/github.com/op/go-logging" ]
then	
	echo 'getting github.com/op/go-logging'
	go get  github.com/op/go-logging
fi
go build main && mv main ./bin