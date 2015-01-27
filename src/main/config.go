package main

type Configuration struct {
	logMainPath,logSubPath,logAccessFileName string
}

var (
	Config = Configuration{
		logMainPath : "/var/log/",
		logSubPath   : "llmf-proxyserver/",
		logAccessFileName :  "access.log",
	}
)

