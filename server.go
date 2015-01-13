package main
import (
	"fmt"
	"io"
	"regexp"
	//"io/ioutil"	
	"net/http"
	"github.com/amahi/spdy"
	"github.com/op/go-logging"
)
var log = logging.MustGetLogger("server")

var httpClient *http.Client = &http.Client{}
var rePath = regexp.MustCompile("^https?://([a-zA-Z0-9\\.\\-]+(\\:\\d+)?)/");

func handler(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("%v", r.Header)
	method := r.Method
	host := r.Header.Get(":host")
	path := r.Header.Get(":path")
	uri:="http://"+host + path
	if m := rePath. FindStringSubmatch(path); m != nil{
		if m[1]==host {
			uri=path
		}
	}
	

	if r.Method == "CONNECT" {
		
	} else {
		req, err := http.NewRequest(method, uri , nil) 
		handleError(err)
		resp, err := httpClient.Do(req)
		handleError(err)
		defer resp.Body.Close()

		for headerKey := range resp.Header {
			headerVal := resp.Header.Get(headerKey)
			w.Header().Set(headerKey, headerVal)
		}
		//body, _ := ioutil.ReadAll(resp.Body)
		io.Copy(w, resp.Body)
	}

}
func handleError(err error) {
	if err != nil {
		//panic(err)
	}
}
func main() {
	//spdy.EnableDebug()
	http.HandleFunc("/", handler)
	log.Info("Launching SPDY on :8080")
	err := spdy.ListenAndServeTLS(":8080", "./cert/serverTLS/server.pem", "./cert/serverTLS/server.key" , nil)
	
	if err != nil {
		fmt.Println(err)
	}else{

	}
	
}