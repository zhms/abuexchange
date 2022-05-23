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

func FuturesTrade() {
	for i := 0; i < len(Symbols); i++ {
		go futures_trade(Symbols[i])
	}
}

func futures_trade(symbol string){
	defer recover()
	conn, _, err := websocket.DefaultDialer.Dial(FututesStreamUrl, nil)
	if err != nil {
		fmt.Println("futures trade:",err)
		time.Sleep(time.Second * 1)
		go market_trade(symbol)
		return
	}
	submsg := SubScribeData{}
	submsg.Id = 1
	submsg.Method = "SUBSCRIBE"
	submsg.Params = []string{}
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol, "/", "", -1), "@trade"))
	sendmsg, _ := json.Marshal(&submsg)
	conn.WriteMessage(1, sendmsg)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		go futures_trade_message(message)
	}
	time.Sleep(time.Second * 1)
	go market_trade(symbol)
}

func futures_trade_message(s []byte) {
	defer recover()
	parser := fastjson.Parser{}
	data,_ := parser.ParseBytes(s)
	stream := string(data.GetStringBytes("stream"))
	if len(stream) == 0 { return }
	symbol := strings.Split(stream,"@")[0]
	tdata := markettradedata{}
	tdata.TradeAmount,_ =  strconv.ParseFloat(string(data.GetStringBytes("data","q")),64)
	tdata.TradePrice,_ =  strconv.ParseFloat(string(data.GetStringBytes("data","p")),64)
	tdata.TradeTime = data.GetInt64("data","E") / 1000
	strdata,_ := json.Marshal(&tdata)
	redisconn := RedisPool.Get()
	redisconn.Do("publish","reptile_futures_trade",fmt.Sprint(symbol,"@",string(strdata)))
	redisconn.Do("publish","reptile_futures_price",fmt.Sprint(symbol,"@",tdata.TradePrice))
	redisconn.Close()
}
