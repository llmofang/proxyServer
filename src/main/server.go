package main
import (
	"myproxy"
	"fmt"
	"net/http"
	"github.com/amahi/spdy"
)
var proxy myproxy.Handler
func main() {
	fmt.Println("%v",proxy)
	spdy.EnableDebug()
	http.HandleFunc("/",proxy.handle)
	//log.Info("Launching SPDY on :8080")
	//fmt.Println("%v",proxy)
	err := spdy.ListenAndServeTLS(":8080", "../cert/serverTLS/server.pem", "../cert/serverTLS/server.key" , nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
