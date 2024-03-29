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

func (h *Helper) GetDataInfo(token string)(int,string,string,int64,string) {//flag,userid,appid,remain,whitelist 
	if token == ""{
		return -1,"0" ,"0",int64(0),""
	}
	data_request, err := h.Client.Hmget("request_token_"+token,"user_id","app_id","whitelist_pattern")
	if err!=nil{
		return -2,"0" ,"0",int64(0),""
	}
	data_data, err := h.Client.Hmget("user_data_"+string(data_request[0])+"_"+string(data_request[1]),"user_id","data_left")
	if err!=nil{
		return -3,"0" ,"0",int64(0),""
	}	
	remain_data,err := strconv.ParseInt(string(data_data[1]),10,64)
	if err != nil{
		return -4,"0" ,"0",int64(0),""
	}
	return 0,string(data_request[0]) ,string(data_request[1]),remain_data,string(data_request[2]) 
}

func (h *Helper) SetDataRemain(key string,remain int64){
	h.Client.Hset("user_data_"+key, "data_left",[]byte(strconv.FormatInt(remain,10)))
}

func (h *Helper) AddData() int64{
	return int64(0)
}

func (h *Helper) MinusData(i int64) int64{
	return int64(0)
}

func (h *Helper) Test(){
	//h.Client.Hmset("request_token_anbo1v1y5",map[string]interface{}{"user_id":"demouser_anbo1v1y5"} )
	//h.Client.Hmset("user_data_demouser_anbo1v1y5",map[string]interface{}{"user_id":"demouser_anbo1v1y5","data_left":8192000,"data_type":0} )
}