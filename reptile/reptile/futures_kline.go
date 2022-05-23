package reptile

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/valyala/fastjson"
)

func FuturesKLine() {
	for i := 0; i < len(Symbols); i++ {
		go futures_kline(Symbols[i])
	}
}

func futures_kline(symbol string){
	defer recover()
	conn, _, err := websocket.DefaultDialer.Dial(FututesStreamUrl, nil)
	if err != nil {
		fmt.Println("futures kline:",err)
		time.Sleep(time.Second * 1)
		go futures_kline(symbol)
		return
	}
	submsg := SubScribeData{}
	submsg.Id = 1
	submsg.Method = "SUBSCRIBE"
	submsg.Params = []string{}
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol,"/","",-1),"@kline_1m"))
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol,"/","",-1),"@kline_3m"))
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol,"/","",-1),"@kline_5m"))
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol,"/","",-1),"@kline_15m"))
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol,"/","",-1),"@kline_30m"))
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol,"/","",-1),"@kline_2h"))
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol,"/","",-1),"@kline_4h"))
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol,"/","",-1),"@kline_6h"))
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol,"/","",-1),"@kline_8h"))
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol,"/","",-1),"@kline_12h"))
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol,"/","",-1),"@kline_1d"))
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol,"/","",-1),"@kline_3d"))
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol,"/","",-1),"@kline_1w"))
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol,"/","",-1),"@kline_1M"))
	sendmsg,_ := json.Marshal(&submsg)
	conn.WriteMessage(1, sendmsg)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		go futures_kline_message(message)
	}
	time.Sleep(time.Second * 1)
	go futures_kline(symbol)
}


type futuresklinedata struct{
	Id int64 `json:"id"`
	Open float64 `json:"open"`
	Close float64 `json:"close"`
	High float64 `json:"high"`
	Low float64 `json:"low"`
	Vol float64 `json:"vol"`
	Count int64 `json:"count"`
	Amount float64 `json:"amount"`
}

func futures_kline_message(s []byte) {
	defer recover()
	parser := fastjson.Parser{}
	data,_ := parser.ParseBytes(s)
	stream := string(data.GetStringBytes("stream"))
	if len(stream) == 0 { return }
	symbol := strings.Split(stream,"@")[0]
	interval := string(data.GetStringBytes("data","k","i"))
	kdata := futuresklinedata{}
	kdata.Id = data.GetInt64("data","k","t") / 1000
	kdata.Open,_ = strconv.ParseFloat(string(data.GetStringBytes("data","k","o")),64)
	kdata.Close,_ = strconv.ParseFloat(string(data.GetStringBytes("data","k","c")),64)
	kdata.High,_ = strconv.ParseFloat(string(data.GetStringBytes("data","k","h")),64)
	kdata.Low,_ = strconv.ParseFloat(string(data.GetStringBytes("data","k","l")),64)
	kdata.Vol,_ = strconv.ParseFloat(string(data.GetStringBytes("data","k","v")),64)
	kdata.Amount,_ = strconv.ParseFloat(string(data.GetStringBytes("data","k","q")),64)
	kdata.Count = data.GetInt64("data","k","n")
	strdata,_ := json.Marshal(&kdata)
	rediskey := fmt.Sprint("reptile:futures:kline:",symbol,":",interval)
	redisconn := RedisPool.Get()
	redisconn.Do("hset",rediskey,kdata.Id,strdata)
	redisconn.Do("publish","reptile_futures_kline",fmt.Sprint(symbol,"@",interval,"@",string(strdata)))
	redisconn.Close()
}


