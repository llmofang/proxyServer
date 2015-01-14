package gossl

import (
    "net"
    "net/http"
    "io"
    "fmt"
    "github.com/amahi/spdy"
)

type Proxy struct {
    pipeConns []net.Conn
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
    if ProxyHeaderFilter(conn, request) {
        //inject something, so close the connection
        defer conn.Close()
        return
    }

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

    me.pipeConns = append(me.pipeConns, conn)
}

func (me *Proxy) KillConnections() {
    for idx := range me.pipeConns {
        me.pipeConns[idx].Close()
    }
    me.pipeConns = make([]net.Conn, 0)
}

func (me *Proxy) Start(port int) {
    go startProxy(me, port)
}

func startProxy(proxy *Proxy, port int) {
    portStr := fmt.Sprintf(":%d", port)
    http.ListenAndServe(portStr, proxy)
}

