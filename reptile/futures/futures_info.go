package futures

import (
	"encoding/json"
	"exreptile/reptile"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/valyala/fastjson"
)

func FuturesInfo() {
	for i := 0; i < len(reptile.Symbols); i++ {
		go futures_info(reptile.Symbols[i])
	}
}

func futures_info(symbol string) {
	defer recover()
	conn, _, err := websocket.DefaultDialer.Dial(reptile.FututesStreamUrl, nil)
	if err != nil {
		fmt.Println("futures info:", err)
		time.Sleep(time.Second * 1)
		go futures_info(symbol)
		return
	}
	submsg := reptile.SubScribeData{}
	submsg.Id = 1
	submsg.Method = "SUBSCRIBE"
	submsg.Params = []string{}
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol, "/", "", -1), "@markPrice@1s"))
	sendmsg, _ := json.Marshal(&submsg)
	conn.WriteMessage(1, sendmsg)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		go futures_info_message(message)
	}
	time.Sleep(time.Second * 1)
	go futures_info(symbol)
}

type futuresinfodata struct {
	MarkPrice   float64
	FundingRate float64
	FundingTime int64
}

func futures_info_message(s []byte) {
	defer recover()
	parser := fastjson.Parser{}
	data, _ := parser.ParseBytes(s)
	stream := string(data.GetStringBytes("stream"))
	if len(stream) == 0 {
		return
	}
	symbol := strings.Split(stream, "@")[0]
	tdata := futuresinfodata{}
	tdata.MarkPrice, _ = strconv.ParseFloat(string(data.GetStringBytes("data", "p")), 64)
	tdata.FundingRate, _ = strconv.ParseFloat(string(data.GetStringBytes("data", "r")), 64)
	tdata.FundingTime = data.GetInt64("data", "T") / 1000
	strdata, _ := json.Marshal(&tdata)
	redisconn := reptile.RedisPool.Get()
	rediskey := fmt.Sprint("reptile:futures:info:", symbol)
	redisconn.Do("set", rediskey, strdata)
	redisconn.Do("publish", "reptile_futures_info", fmt.Sprint(symbol, "@", string(strdata)))
	redisconn.Close()
}
