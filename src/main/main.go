package main

import(
	"myproxy"
	"log"
	"net/http"
)

func main(){
	file := &File{}
	file.mkdir(Config.logMainPath+Config.logSubPath,750)
	file.cd(Config.logMainPath+Config.logSubPath)
	accesslog :=file.touch(Config.logAccessFileName,610)
	handler := myproxy.NewHandler(accesslog)
	err := http.ListenAndServe(":8080",handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}