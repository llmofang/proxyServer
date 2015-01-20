package myproxy

import(
	"net"
	"net/http"
)

type Handler struct {
	pipeConns []net.Conn
}



func handleError(err error) {
	if err != nil {
		//panic(err)
	}
}

func  (h *Handler) handle(w http.ResponseWriter, r *http.Request) {
	log.Info("NEW REQUEST")
	if r.Method == "CONNECT" {
		h.proxyHttps(w,r)
	} else {
		//h.proxyHttp(w,r)
	}
}