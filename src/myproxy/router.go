package myproxy

import(
	"radius"
	"bytes"
	//"fmt"
	"io"
	"net"
	"net/http"
)

var rh radius.Helper

type Handler struct {
	logFile string
}

type buffer struct {
	bytes.Buffer
}

type user struct {
	authorization string
	total int64
	remain int64
}

func (b *buffer) Close() error {
	return nil
}

func NewHandler(logfile string,redisaddr string) *Handler {
	rh =radius.Helper{
		Addr :redisaddr,
	}
	rh.Init();
	rh.Test();
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

func pipeAndCount(conn1 net.Conn, conn2 net.Conn)  int64{
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

func  (h *Handler) proxyHttp(w http.ResponseWriter, r *http.Request,u *user){
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	transport := &http.Transport{}
	buf := new(buffer)
	requestURL := r.URL.String()
	io.Copy(buf, r.Body)
	newRequest, err := http.NewRequest(r.Method, requestURL, buf)
	//handleError(err)
	copyHeader(r.Header,newRequest.Header)
	l :=int64(len(buf.String()))
	if l > u.remain{
		rh.SetDataRemain(u.authorization,0)
		http.Error(w,"Unauthorized", http.StatusUnauthorized)
		return	
	}
	u.remain -= l
	rh.SetDataRemain(u.authorization,u.remain)
	//u.remain = 
	newRequest.ContentLength = l
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
	l = newResponse.ContentLength
	if l==-1{
		io.Copy(buf,newResponse.Body)
		//log.Error(buf.String())
		l = int64(len(buf.String()))
		newResponse.ContentLength = l
		if l > u.remain{
			rh.SetDataRemain(u.authorization,0)
			http.Error(w,"Unauthorized", http.StatusUnauthorized)
			return	
		}
		copyHeader(newResponse.Header, w.Header())
		w.WriteHeader(newResponse.StatusCode)
		io.Copy(w,buf)
	}else{
		if l > u.remain{
			
			rh.SetDataRemain(u.authorization,0)
			http.Error(w,"Unauthorized", http.StatusUnauthorized)
			return	
		}
		copyHeader(newResponse.Header, w.Header())
		w.WriteHeader(newResponse.StatusCode)
		io.Copy(w,newResponse.Body)
	}
	buf.Close()
	webLog := &webLogger{
		file : h.logFile,
	} 	
	go func(){
	 	webLog.formatLog(ip,"-",r.Method,requestURL,r.Proto ,newResponse.StatusCode,newResponse.ContentLength, r.Header.Get("User-Agent")) 
	 	webLog.write()
		//webLog.dumpLog()
	}()
}


func  (h *Handler) proxyHttps(w http.ResponseWriter, r *http.Request,u *user) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
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
	total := pipeAndCount(serverConn,conn)
	webLog := &webLogger{
		file : h.logFile,
	} 	
	go func(){
	 	webLog.formatLog(ip,"-","CONNECT","https://"+r.Host,r.Proto ,200,total, r.Header.Get("User-Agent")) 
	 	webLog.write()
		//webLog.dumpLog()
	}()
	//fmt.Println(total)
}



func (h *Handler) DebugURL(response http.ResponseWriter, request *http.Request) {

}

func  (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "" {
		http.Error(w,"Bad Request", http.StatusBadRequest)
		return
	}
	authorization:= r.Header.Get("Llmf-Proxy-Authorization")
	total,remain :=  rh.GetDataInfo(authorization) 
	/*
	if authorization == "" || -1 == remain{
		http.Error(w,"Unauthorized", http.StatusUnauthorized)
		return
	}
	*/
	u := &user{
		authorization:authorization,
		total:total,
		remain:remain,
	}
	if r.Method == "CONNECT" {
		h.proxyHttps(w,r,u)
	} else {
		h.proxyHttp(w,r,u)
	}
}