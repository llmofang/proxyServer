package main

import(
	"myproxy"
	//"fmt"
	//"log"
	"net/http"
	"github.com/amahi/spdy"
)

func startNormalProxy(h http.Handler){
	http.ListenAndServe(Config.port,h)
}

func startSpdyProxy(h http.Handler){
	spdy.ListenAndServeTLS(":8081", "../cert/serverTLS/server.pem", "../cert/serverTLS/server.key",h)
}

func main(){
	command:= &Command{}
	command.mkdir(Config.logMainPath+Config.logSubPath,750)
	command.cd(Config.logMainPath+Config.logSubPath)
	accesslog :=command.touch(Config.logAccessFileName,610)
	handler := myproxy.NewHandler(accesslog,Config.redisServerAddr)
	go startSpdyProxy(handler)
	startNormalProxy(handler)
}