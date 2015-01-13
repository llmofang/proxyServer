package main
import (
	"fmt"
	"io"
	"bufio"
	"regexp"
	//"io/ioutil"
	"net"	
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
	scheme:=r.Header.Get(":scheme")
	uri:=scheme+host+path

	if m := rePath. FindStringSubmatch(path); m != nil{
		if m[1]==host {
			uri=path
		}
	}
	
	fmt.Printf("%v", r.Header)
	log.Error(r.Method)
	if r.Method == CONNECT {
		conn, err := net.Dial("tcp",host)
		if err != nil {
			log.Printf("Conection failed: %v", err)
			return
		}
		fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
		status, err := bufio.NewReader(conn).ReadString('\n')

		fmt.Printf("%v", status)


	} else {
		req, err := http.NewRequest(method, uri , nil) 
		handleError(err)
		

		for headerKey := range r.Header{
			headerVal := r.Header.Get(headerKey)

			req.Header.Set(headerKey, headerVal)
		}

		req.Header.Set("User-Agent", "LLMF PROXY SERVER")

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
	spdy.EnableDebug()
	http.HandleFunc("/", handler)
	log.Info("Launching SPDY on :8080")
	err := spdy.ListenAndServeTLS(":8080", "./cert/serverTLS/server.pem", "./cert/serverTLS/server.key" , nil)
	
	if err != nil {
		fmt.Println(err)
	}else{

	}
	
}