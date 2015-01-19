package main
import(
	"flag"
	"fmt"
	"net"
)
var connid = uint64(0)
var Addr = flag.String("l", ":9999", "local address")
func handle(err error) {
	if err != nil {
		panic(err)
	}
}
func main(){
	addr, err := net.ResolveTCPAddr("tcp", *Addr)	
	handle(err);
	listener, err := net.ListenTCP("tcp", addr)
	handle(err);
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("Failed to accept connection '%s'\n", err)
			continue
		}
		fmt.Printf("%v",conn)
		connid++
		fmt.Printf("%v",connid)
	}	

}