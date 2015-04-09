package myproxy

import(
	"radius"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
)

var rh radius.Helper

type Handler struct {
	logFile string
}

type buffer struct {
	bytes.Buffer
}

type user struct {
	userid string
	appid string
	remain int64
}

func (b *buffer) Close() error {
	return nil
}

func NewHandler(logfile string,redisaddr string) *Handler {

	client := radius. InitClient(redisaddr,0,"",100)
	rh =radius.Helper{
		Client:client,
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

func spdyHeaderToNormal(h *http.Header){
	var hArr [] string = [] string{":method",":host",":path",":version",":scheme",":Via"}
	for _, v:= range hArr{
		h.Del(v)
	}
}

func  (h *Handler) proxyHttp(w http.ResponseWriter, r *http.Request,u *user){
	var requestURL string
	var data_key string
	data_key  = u.userid+"_"+u.appid
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	l := int64(0)
	buf := new(buffer)
	if r.Header.Get(":method")!=""{
		requestURL = "http://"+ r.Header.Get(":host") + r.Header.Get(":path")
		spdyHeaderToNormal(&r.Header)
	}else{
		requestURL = r.URL.String()
	}
	transport := &http.Transport{}
	if r.Method == "POST" || r.Method =="PUT" {
		io.Copy(buf, r.Body)// to fix
		l :=int64(len(buf.String()))
		if l > u.remain{
			rh.SetDataRemain(data_key,0)
			h.ServeError(w,r,"100","FLOW EXHAUSTED WHEN POST DATA")
			return	
		}
		u.remain -= l	
		rh.SetDataRemain(data_key,u.remain)	
	}
	newRequest, err := http.NewRequest(r.Method, requestURL, buf)
	//handleError(err)
	buf.Reset()
	copyHeader(r.Header,newRequest.Header)
	newRequest.ContentLength = l
	newResponse, err := transport.RoundTrip(newRequest)
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
		l = int64(len(buf.Bytes()))//int64(buf.Len())
		newResponse.ContentLength = l
		if l > u.remain{
			rh.SetDataRemain(data_key,0)
			h.ServeError(w,r,"100","FLOW EXHAUSTED WHEN GET CONTENT")
			return	
		}
		copyHeader(newResponse.Header, w.Header())
		w.WriteHeader(newResponse.StatusCode)
		io.Copy(w,buf)
	}else{
		if l > u.remain{
			rh.SetDataRemain(data_key,0)
			h.ServeError(w,r,"100","FLOW EXHAUSTED WHEN GET CONTENT")
			return	
		}
		copyHeader(newResponse.Header, w.Header())
		w.WriteHeader(newResponse.StatusCode)
		io.Copy(w,newResponse.Body)
	}
	buf.Close()
	u.remain -=l
	rh.SetDataRemain(data_key,u.remain)
	webLog := &webLogger{
		file : h.logFile,
	} 		
	go func(){
	 	webLog.formatLog(ip,u.userid,u.appid,r.Method,requestURL,r.Proto ,newResponse.StatusCode,newResponse.ContentLength, r.Header.Get("User-Agent")) 
	 	webLog.write()
		//webLog.dumpLog()
	}()
}


func  (h *Handler) proxyHttps(w http.ResponseWriter, r *http.Request,u *user) {
	var data_key string
	data_key  = u.userid+"_"+u.appid
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
	rh.SetDataRemain(data_key,u.remain)
	
	defer conn.Close()
	defer serverConn.Close()

	if total > u.remain{
		u.remain -= 0	
	}
	u.remain -= total	
	rh.SetDataRemain(data_key,u.remain)	
	webLog := &webLogger{
		file : h.logFile,
	} 	
	go func(){
	 	webLog.formatLog(ip,u.userid,u.appid,"CONNECT","https://"+r.Host,r.Proto ,200,total, r.Header.Get("User-Agent")) 
	 	webLog.write()
		//webLog.dumpLog()
	}()
	//fmt.Println(total)
}

func (h *Handler) DebugURL(response http.ResponseWriter, request *http.Request) {

}

func (h *Handler) ServeError(w http.ResponseWriter, r *http.Request,code string,message string){
	if r.Method == "CONNECT"{
		conn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			conn.Close()
		}		
		conn.Write([]byte("HTTP/1.1 401 Authorization Required\r\n" +
			"Content-Type: text/html\r\n" +
			"Content-Length: 200\r\n" +
			"\r\n"));
		defer conn.Close()
	}else{
		fmt.Println("%v",message)
		w.Header().Set("Content-Type","application/json")
		http.Error(w,"{\"code\":"+code+" ,\" message\":"+message+"} ", http.StatusUnauthorized)
	}
	
}

func  (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "" {
		http.Error(w,"Bad Request", http.StatusBadRequest)
		return
	}
	authorization:= r.Header.Get("Proxy-Authorization")
	fmt.Println("%v",r.Header)
	fmt.Println("%v",r.Method)
	//fmt.Println("%v",authorization)
	log.Info(authorization);
	//authorization = "f72e903dab05735ad7d15008193f84b1c6f9a2d01928fc6fbe760bd47eadc8fb4e577f06be430b9743239d64a2a7b8cba060141a7c5cc4464d02c6daa80f275e"
	r.Header.Del("Proxy-Authorization")
	flag,userid,appid,remain,whitelist :=  rh.GetDataInfo(authorization) 
	if 0 != flag {
		var errMsg string 
		switch flag{
			case -1:
				errMsg = "EMPTY TOKEN";
			break;
			case -2:
				errMsg = "TOKEN EXPIRED";
			break;
			case -3:
				errMsg = "NO AVAILABLE FLOW";
			break;
			case -4:
				errMsg = "NO FLOW";
			break;
		}
		h.ServeError(w,r,"100",errMsg)
		return
	}
	whitelist = ""
	if "" != whitelist{
		match, _ := regexp.MatchString(whitelist,r.URL.String())
		fmt.Println(match)
		if !match{
			h.ServeError(w,r,"101","URL NOT MATCH WHITELIST")
			return
		}
	}
	u := &user{
		userid:userid,
		appid:appid,
		remain:remain,
	}
	if r.Method == "CONNECT" {
		h.proxyHttps(w,r,u)
	} else {
		h.proxyHttp(w,r,u)
	}
}