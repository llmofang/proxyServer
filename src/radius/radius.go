package radius
import(
	//"fmt"
	"github.com/hoisie/redis"
)

type Helper struct {
	Addr string
	client redis.Client
}

func (h *Helper) Init(){
	h.client.Addr = h.Addr
}

func (h *Helper) GetDataInfo(token string)(int64,int64) {
	//return int64(1024000) ,int64(30000)
}

func (h *Helper) SetDataRemain(token string,remain int64){

}

func (h *Helper) AddData() int64{
	return int64(0)
}

func (h *Helper) MinusData(i int64) int64{
	return int64(0)
}



func (h *Helper) Test(){

}