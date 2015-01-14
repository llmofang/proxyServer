package main
import (
	//"gossl"
	"fmt"
	//"io"
	//"os"
	//"bufio"
	"regexp"
	//"io/ioutil"
	//"net"	
	"net/http"
	"github.com/amahi/spdy"
)

type Proxy struct {

}

func (me *Proxy) ServeHTTP(response http.ResponseWriter, request *http.Request) {
    fmt.Println(`request`);
    if request.Method == "CONNECT" {
        me.proxyHttps(response, request)
    } else {
        serverRequest(false, response, request)
    }
}

func (me *Proxy) proxyHttps(response http.ResponseWriter, request *http.Request) {

}



func  serverRequest(isSSL bool, response http.ResponseWriter, request *http.Request){

}

var rePath = regexp.MustCompile("^https?://([a-zA-Z0-9\\.\\-]+(\\:\\d+)?)/");
/*
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v",r.Header)

	return
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

	
	if r.Method == "CONNECT" {
		//local, err := net.Listen("tcp", localAddr)
		//conn, err := net.Dial("tcp",host)
		//if err != nil {
		//	//log.Printf("Conection failed: %v", err)
		//	return
		//}
		//fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
		//status, err := bufio.NewReader(conn).ReadString('\n')
//
		//fmt.Printf("%v", status)


	} else {
		/*
		req, err := http.NewRequest(method, uri , nil) 
		handleError(err)
		log.Info("Feach: " + uri)

		for headerKey := range r.Header{
			headerVal := r.Header.Get(headerKey)
			req.Header.Set(headerKey, headerVal)
		}

		//req.Header.Set("User-Agent", "LLMF PROXY SERVER")

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
*/

func handleError(err error) {
	if err != nil {
		//panic(err)
	}
}
func startProxy(proxy *Proxy) {
	err := spdy.ListenAndServeTLS(":8080", "./cert/serverTLS/server.pem", "./cert/serverTLS/server.key" , proxy)

	if err != nil {
		fmt.Println(err)
		//log.Error(err)
	}
}

func main() {
	//spdy.EnableDebug()
	//http.HandleFunc("/", proxy.ServeHTTP)
	//log.Info("Launching SPDY on :8080")
	proxy:= new(Proxy)
	go startProxy(proxy)




}
