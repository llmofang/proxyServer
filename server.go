package main
import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"github.com/SlyMarbo/spdy"
    	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("example")

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func handleProxy(conn spdy.Conn) {
	
	url := "http://" + conn.Conn().RemoteAddr().String() + "/"
	log.Debug("debug %s", string(url))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	res, err := conn.RequestResponse(req, nil, 1)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		panic(err)
	}

	res.Body.Close()

	fmt.Println(buf.String())
}

func main() {
	handler := spdy.ProxyConnHandlerFunc(handleProxy)
	http.Handle("/", spdy.ProxyConnections(handler))
	handle(http.ListenAndServeTLS(":8080", "./cert/serverTLS/server.pem", "./cert/serverTLS/server.key", nil))
	//handle(http.ListenAndServe(":8080", nil))
	log.Info("Launching SPDY on :8080")
}
