package myproxy
import "strconv"
func formatLog(ip string,user string,time string,method string,path string,proto string ,status int, size int64, agent string) string{
	return ip+" - "+user+" ["+time+"] \""+method+" "+path+" "+proto+"\" "+strconv.Itoa(status)+" "+strconv.FormatInt(size,10)+" \"-\" \""+agent+"\""
}

