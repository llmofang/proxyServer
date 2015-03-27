package main

import(
	"myproxy"
	//"fmt"
	//"log"
	"net/http"
	//"github.com/amahi/spdy"
)

func startNormalProxy(h http.Handler){
	http.ListenAndServe(Config.port,h)
}

func startHttp2Proxy(h http.Handler){
}

func prepareLogFile() string{
	command:= &Command{}
	command.mkdir(Config.logMainPath+Config.logSubPath,750)
	command.cd(Config.logMainPath+Config.logSubPath)
	return command.touch(Config.logAccessFileName,610)	
}
func main(){
	accesslog := prepareLogFile();
	handler := myproxy.NewHandler(accesslog,Config.redisServerAddr)
	startNormalProxy(handler)
}