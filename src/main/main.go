package main

import(
	"myproxy"
	"log"
	"net/http"
)
func main(){
	command:= &Command{}
	command.mkdir(Config.logMainPath+Config.logSubPath,750)
	command.cd(Config.logMainPath+Config.logSubPath)
	accesslog :=command.touch(Config.logAccessFileName,610)
	handler := myproxy.NewHandler(accesslog,Config.redisServerAddr)
	log.Println("Starting Server At ", Config.port)
	err := http.ListenAndServe(Config.port,handler)
	if err != nil {
		log.Fatal("Start Error: ", err)
	}	
}