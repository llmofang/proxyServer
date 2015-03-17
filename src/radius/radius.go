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

func (h *Helper) GetDataInfo(token string)(string,int,int64,string) {//userid,appid,remain,whitelist 
	if token == ""{
		//return "0" ,int64(-1)
	}
	data_request, err := h.Client.Hmget("request_token_"+token,"user_id","app_id")
	if err!=nil{
		//return "0" ,int64(-1)
	}
	data_data, err := h.Client.Hmget("user_data_"+string(data_request[0])+"_"+string(data_request[1]),"user_id","data_left","whitelist_pattern")
	if err!=nil{
		//return "0" ,int64(-1)
	}	
	remain_data,err := strconv.ParseInt(string(data_data[1]),10,64)
	if err != nil{
		r = -1
	}
	return string(data_request[0]) ,int(data_request[1]),remain_data,string(data_data[2]) 
}

func (h *Helper) SetDataRemain(token string,remain int64){
	//fmt.Println("%v\n",remain)
	//h.Client.Hset("request_token_"+token, "data_left",[]byte(strconv.FormatInt(remain,10)))
}

func (h *Helper) AddData() int64{
	return int64(0)
}

func (h *Helper) MinusData(i int64) int64{
	return int64(0)
}

func (h *Helper) Test(){
	h.Client.Hmset("request_token_anbo1v1y5",map[string]interface{}{"user_id":"demouser_anbo1v1y5"} )
	h.Client.Hmset("user_data_demouser_anbo1v1y5",map[string]interface{}{"user_id":"demouser_anbo1v1y5","data_left":8192000,"data_type":0} )
}