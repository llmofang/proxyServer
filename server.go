package main
import (
	"fmt"
	"net/http"
	"github.com/amahi/spdy"
	"github.com/op/go-logging"
)
var log = logging.MustGetLogger("server")
func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "CONNECT" {
		
	} else {
		log.Debug(r.Header.Get("host"))
		//http.NewRequest(r.Method, urlStr string, body io.Reader) (*Request, error)
	}
}
func main() {
	spdy.EnableDebug()
	http.HandleFunc("/", handler)
	log.Info("Launching SPDY on :8080")
	err := spdy.ListenAndServeTLS(":8080", "./cert/serverTLS/server.pem", "./cert/serverTLS/server.key" , nil)

	if err != nil {
		fmt.Println(err)
	}else{

	}
	
}