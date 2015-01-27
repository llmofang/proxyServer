package myproxy

import(
	"bytes"
	//"fmt"
	"time"
	"io"
	"net"
	"regexp"
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

func writeLog(){

}

func  copyHeader(from, to http.Header) {
	for hdr, items := range from {
		for _, item := range items {
			to.Add(hdr, item)
		}
	}
}
func parseRequestUrl(w *http.Request) string {
	var requestURL string
	if m, _ := regexp.MatchString("^http[s]{0,1}://.*$", w.RequestURI); m {
		requestURL = w.RequestURI
	} else {
		requestURL = "http://" + w.Host + w.RequestURI
	}
	return requestURL
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
	newResponse, err := transport.RoundTrip(newRequest)
	//fmt.Println("%v",newResponse)
	if err != nil {
		http.NotFound(w,r)
		return
	}	
	defer newResponse.Body.Close()

	//_, err = io.Copy(buf, newResponse.Body)

	//handleError(err)

	//fmt.Fprintf(response,  buf.String())



	copyHeader(newResponse.Header, w.Header())

	

	if newResponse.StatusCode!=200{
		
	}

	w.WriteHeader(newResponse.StatusCode)

	io.Copy(w,newResponse.Body)

	logstr := formatLog(ip,"-",time.Now().Format("02/Jan/2006:15:04:05 -0700"),r.Method,requestURL,r.Proto ,newResponse.StatusCode,newResponse.ContentLength, r.Header.Get("User-Agent")) 
	log.Info(logstr)
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

	go io.Copy(serverConn, conn)
	go io.Copy(conn, serverConn)

	//h.pipeConns = append(h.pipeConns, conn)
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