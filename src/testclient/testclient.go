package main
import(
	"flag"
	"net"
)
var Addr = flag.String("l", ":9999", "local address")
func handle(err error) {
	if err != nil {
		panic(err)
	}
}
func main(){
	addr, err := net.ResolveTCPAddr("tcp", *Addr)
	handle(err);
	net.DialTCP("tcp", nil, addr)	
	
}