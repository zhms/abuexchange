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

func FuturesTicker() {
	for i := 0; i < len(reptile.Symbols); i++ {
		go futures_ticker(reptile.Symbols[i])
	}
}

func futures_ticker(symbol string) {
	defer recover()
	conn, _, err := websocket.DefaultDialer.Dial(reptile.FututesStreamUrl, nil)
	if err != nil {
		fmt.Println("futures ticker:", err)
		time.Sleep(time.Second * 1)
		go futures_ticker(symbol)
		return
	}
	submsg := reptile.SubScribeData{}
	submsg.Id = 1
	submsg.Method = "SUBSCRIBE"
	submsg.Params = []string{}
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol, "/", "", -1), "@ticker"))
	sendmsg, _ := json.Marshal(&submsg)
	conn.WriteMessage(1, sendmsg)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		go futures_ticker_message(message)
	}
	time.Sleep(time.Second * 1)
	go futures_ticker(symbol)
}

type futurestickerdata struct {
	Id     int64   `json:"id"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Vol    float64 `json:"vol"`
	Count  int64   `json:"count"`
	Amount float64 `json:"amount"`
	Change float64 `json:"change"`
}

func futures_ticker_message(s []byte) {
	defer recover()
	parser := fastjson.Parser{}
	data, _ := parser.ParseBytes(s)
	stream := string(data.GetStringBytes("stream"))
	if len(stream) == 0 {
		return
	}
	symbol := strings.Split(stream, "@")[0]
	tdata := futurestickerdata{}
	tdata.Id = data.GetInt64("data", "E") / 1000
	tdata.Open, _ = strconv.ParseFloat(string(data.GetStringBytes("data", "o")), 64)
	tdata.Close, _ = strconv.ParseFloat(string(data.GetStringBytes("data", "c")), 64)
	tdata.High, _ = strconv.ParseFloat(string(data.GetStringBytes("data", "h")), 64)
	tdata.Low, _ = strconv.ParseFloat(string(data.GetStringBytes("data", "l")), 64)
	tdata.Vol, _ = strconv.ParseFloat(string(data.GetStringBytes("data", "v")), 64)
	tdata.Amount, _ = strconv.ParseFloat(string(data.GetStringBytes("data", "q")), 64)
	tdata.Change, _ = strconv.ParseFloat(string(data.GetStringBytes("data", "P")), 64)
	tdata.Count = data.GetInt64("data", "n")
	strdata, _ := json.Marshal(&tdata)
	rediskey := fmt.Sprint("reptile:futures:ticker:", symbol)
	redisconn := reptile.RedisPool.Get()
	redisconn.Do("set", rediskey, strdata)
	redisconn.Do("publish", "reptile_futures_ticker", fmt.Sprint(symbol, "@", string(strdata)))
	redisconn.Close()
}
