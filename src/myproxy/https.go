package myproxy
import(
	"io"
	"net"
	"fmt"
	"net/http"
)
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
