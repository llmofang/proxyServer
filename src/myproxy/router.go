package myproxy

import(
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
)

type Handler struct {
	logFile string
	pipeConns []net.Conn
}

type buffer struct {
	bytes.Buffer
}

func (b *buffer) Close() error {
	return nil
}

func NewHandler(logfile string) *Handler {
	return &Handler{
		logFile : logfile,
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}


func cacheFile(){

}


func  copyHeader(from, to http.Header) {
	for hdr, items := range from {
		for _, item := range items {
			to.Add(hdr, item)
		}
	}
}

func chanFromConn(conn net.Conn) chan []byte {
	c := make(chan []byte)
	go func() {
		b := make([]byte, 1024)
		for {
			n, err := conn.Read(b)
			if n > 0 {
				res := make([]byte, n)
				// Copy the buffer so it doesn't get changed while read by the recipient.
				copy(res, b[:n])
				c <- res
			}
			if err != nil {
				c <- nil
				break
			}
		}
	}()

	return c
}

func Pipe(conn1 net.Conn, conn2 net.Conn)  int64{
	chan1 := chanFromConn(conn1)
	chan2 := chanFromConn(conn2)
	total := 0
	for {
		total ++
		select {
			case b1 := <-chan1:
			if b1 == nil {
				return int64(total*1024) 
			} else {
				conn2.Write(b1)
			}
			case b2 := <-chan2:
			if b2 == nil {
				return int64(total*1024) 
			} else {
				conn1.Write(b2)
			}
		}
	}
}

func  (h *Handler) proxyHttp(w http.ResponseWriter, r *http.Request){
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	transport := &http.Transport{}
	buf := new(buffer)
	requestURL := r.URL.String()
	io.Copy(buf, r.Body)
	newRequest, err := http.NewRequest(r.Method, requestURL, buf)
	//handleError(err)
	copyHeader(r.Header,newRequest.Header)
	l :=len(buf.String())
	newRequest.ContentLength = int64(l)
	newResponse, err := transport.RoundTrip(newRequest)
	buf.Reset()
	//fmt.Println("%v",newResponse)
	if err != nil {
		http.NotFound(w,r)
		return
	}	
	defer newResponse.Body.Close()
	//_, err = io.Copy(buf, newResponse.Body)
	//handleError(err)
	copyHeader(newResponse.Header, w.Header())
	w.WriteHeader(newResponse.StatusCode)
	
	webLog := &webLogger{
		file : h.logFile,
	} 
	if newResponse.ContentLength==-1{
		io.Copy(buf,newResponse.Body)
		//log.Error(buf.String())
		l :=len(buf.String())
		newResponse.ContentLength = int64(l)
		io.Copy(w,buf)
	}else{
		io.Copy(w,newResponse.Body)
	}
	buf.Close()
	go func(){
	 	webLog.formatLog(ip,"-",r.Method,requestURL,r.Proto ,newResponse.StatusCode,newResponse.ContentLength, r.Header.Get("User-Agent")) 
	 	webLog.write()
		//webLog.dumpLog()
	}()
}


func  (h *Handler) proxyHttps(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	conn, _, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	serverConn, err := net.Dial("tcp",r.Host)
	if err != nil {
		return
	}
	conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n" +
	"Content-Type: text/html\r\n" +
	"Content-Length: 200\r\n" +
	"\r\n"));
	// go io.Copy(serverConn,conn)
	// go io.Copy(conn,serverConn)
	total := Pipe(serverConn,conn)
	fmt.Println(total)
}



func (h *Handler) DebugURL(response http.ResponseWriter, request *http.Request) {


}

func  (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//log.Info("NEW REQUEST")
	if r.Method == "" {
		http.Error(w,"Bad Request", http.StatusBadRequest)
		return
	}
	if r.Method == "CONNECT" {
		h.proxyHttps(w,r)
	} else {
		h.proxyHttp(w,r)
	}
}