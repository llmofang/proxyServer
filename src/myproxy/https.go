package myproxy
import(
	"io"
	"net"
	"net/http"
)

type Proxy struct {
    pipeConns []net.Conn
}

func (p *Proxy) proxyHttps(response http.ResponseWriter, request *http.Request) {
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
	serverConn, err := net.Dial("tcp", request.Host)
	if err != nil {
		return
	}

	conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n" +
	"Content-Type: text/html\r\n" +
	"Content-Length: 200\r\n" +
	"\r\n"));

	go io.Copy(serverConn, conn)
	go io.Copy(conn, serverConn)

	p.pipeConns = append(p.pipeConns, conn)
}