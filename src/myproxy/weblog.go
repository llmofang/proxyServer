package myproxy
import(
	//"io"
	"fmt"
	"os"
	"time"
	"strconv"
)
type webLogger struct {
	file string
	log string
}

func (w *webLogger)formatLog(ip string,user string,method string,path string,proto string ,status int, size int64, agent string){
	w.log=ip+" - "+user+" ["+time.Now().Format("02/Jan/2006:15:04:05 -0700")+"] \""+method+" "+path+" "+proto+"\" "+strconv.Itoa(status)+" "+strconv.FormatInt(size,10)+" \"-\" \""+agent+"\""
}
 
func (w *webLogger) dumpLog(){ 
	log.Info(w.log)
}

func (w *webLogger) write(){ 
	f, err := os.OpenFile(w.file,os.O_APPEND|os.O_WRONLY,0600) 
	if err!= nil{
	}
	defer f.Close()
	if _, err = f.WriteString(w.log+"\n"); err != nil {
		fmt.Println("%v",err)
	}
}

