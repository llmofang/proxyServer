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
	handler := myproxy.NewHandler(accesslog)
	err := http.ListenAndServe(":8080",handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}