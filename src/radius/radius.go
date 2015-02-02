package radius
import(
	//"fmt"
	"strconv"
	"github.com/hoisie/redis"
)

type Helper struct {
	Client redis.Client
}

func (h *Helper) Init(){

}

func InitClient(addr string,db int,password string,maxpoolsize int )  redis.Client{
	client := redis.Client{
		Addr:addr,
		//Db:db,
		//Password:password,
		//MaxPoolSize:maxpoolsize,
	}
	return client
}

func (h *Helper) GetDataInfo(token string)(string,int64) {
	if token == ""{
		return "0" ,int64(-1)
	}
	data, err := h.Client.Hmget("request_token_"+token,"user_id","data_left")
	if err!=nil{
		//fmt.Println("%v",err)
		return "0" ,int64(-1)
	}
	r,err := strconv.ParseInt(string(data[1]),10,64)
	if err != nil{
		r = -1
	}
	//fmt.Println("%v",r)
	return string(data[0]) ,r
}

func (h *Helper) SetDataRemain(token string,remain int64){
	//fmt.Println("%v\n",remain)
	h.Client.Hset("request_token_"+token, "data_left",[]byte(strconv.FormatInt(remain,10)))
}

func (h *Helper) AddData() int64{
	return int64(0)
}

func (h *Helper) MinusData(i int64) int64{
	return int64(0)
}



func (h *Helper) Test(){
	h.Client.Hmset("request_token_anbo1v1y5",map[string]interface{}{"user_id":1,"data_left":8192000,"data_type":0} )
}