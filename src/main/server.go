package main
import (
	"fmt"
	"io"
	//"os"
	//"bufio"
	"regexp"
	//"io/ioutil"
	"net"	
	"net/http"
	"github.com/amahi/spdy"
)
var httpClient *http.Client = &http.Client{}
var rePath = regexp.MustCompile("^https?://([a-zA-Z0-9\\.\\-]+(\\:\\d+)?)/");

type Handler struct {
	host string
	pipeConns []net.Conn
}
func NewHandler() *Handler {
	return &Handler{}
}
func  (h *Handler) proxyHttp(response http.ResponseWriter, request *http.Request){

}
func  (h *Handler) proxyHttps(response http.ResponseWriter, request *http.Request) {
	
	fmt.Println("%v",request.Header.Get(":host"))
	hj, ok := response.(http.Hijacker)
	if !ok {
		http.Error(response, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	conn, _, err := hj.Hijack()
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	//check filters
	//if ProxyHeaderFilter(conn, request) {
	//	//inject something, so close the connection
	//	defer conn.Close()
	//	return
	//}

	//process real https proxy
	serverConn, err := net.Dial("tcp", request.Header.Get(":host"))
	if err != nil {
		return
	}

	conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n" +
	"Content-Type: text/html\r\n" +
	"Content-Length: 200\r\n" +
	"\r\n"));

	go io.Copy(serverConn, conn)
	go io.Copy(conn, serverConn)

	h.pipeConns = append(h.pipeConns, conn)
}

func  (h *Handler) handle(w http.ResponseWriter, r *http.Request) {
	log.Info("NEW REQUEST")
	fmt.Println("%v",r)
	if r.Method == "CONNECT" {
		//h.proxyHttps(w,r)
	} else {
		//h.proxyHttp(w,r)
	}
}


func handleError(err error) {
	if err != nil {
		//panic(err)
	}
}

func main() {
	//spdy.EnableDebug()
	handler := NewHandler()
	http.HandleFunc("/", handler.handle)
	//log.Info("Launching SPDY on :8080")
	//fmt.Println("%v",proxy)
	err := spdy.ListenAndServeTLS(":8080", "../cert/serverTLS/server.pem", "../cert/serverTLS/server.key" , nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}



}
