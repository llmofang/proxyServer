package myproxy

import(
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
)

type Handler struct {
	pipeConns []net.Conn
}

type buffer struct {
	bytes.Buffer
}

func (b *buffer) Close() error {
	return nil
}

func NewHandler() *Handler {
	return &Handler{}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func  (h *Handler) proxyHttp(response http.ResponseWriter, request *http.Request){
	fmt.Println("%v",request.Header)

	client := &http.Client{}
	//buf := new(buffer)
	requestURL := request.Header.Get(":scheme") +"://"+ request.Header.Get(":host") + request.Header.Get(":path")
	log.Info(requestURL)
	newRequest, err := http.NewRequest(request.Method, requestURL, nil)
	handleError(err)
	newResponse, err := client.Do(newRequest)
	handleError(err)
	defer newResponse.Body.Close()
	//_, err = io.Copy(buf, newResponse.Body)
	//handleError(err)

	//fmt.Fprintf(response,  buf.String())
	io.Copy(response,newResponse.Body)


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


func (h *Handler) DebugURL(response http.ResponseWriter, request *http.Request) {


}

func  (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	log.Info("NEW REQUEST")
	log.Error(r.Method)
	fmt.Println("%v",r.Header)
	if r.Method == "CONNECT" {
		h.proxyHttps(w,r)
	} else {
		h.proxyHttp(w,r)
	}
}