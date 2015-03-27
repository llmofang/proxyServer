package main

type Configuration struct {
	port,logMainPath,logSubPath,logAccessFileName,redisServerAddr string
}

var (
	Config = Configuration{
		port : ":8080",
		logMainPath : "/var/log/",
		logSubPath   : "llmf-proxyserver/",
		logAccessFileName :  "access.log",
		redisServerAddr : "127.0.0.1:6379",
	}
)

