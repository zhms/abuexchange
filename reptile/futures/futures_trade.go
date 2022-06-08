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

func FuturesTrade() {
	for i := 0; i < len(reptile.Symbols); i++ {
		go futures_trade(reptile.Symbols[i])
	}
}

func futures_trade(symbol string) {
	defer recover()
	conn, _, err := websocket.DefaultDialer.Dial(reptile.FututesStreamUrl, nil)
	if err != nil {
		fmt.Println("futures trade:", err)
		time.Sleep(time.Second * 1)
		go futures_trade(symbol)
		return
	}
	submsg := reptile.SubScribeData{}
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
	go futures_trade(symbol)
}

type futurestradedata struct {
	TradeTime   int64
	TradeAmount float64
	TradePrice  float64
}

func futures_trade_message(s []byte) {
	defer recover()
	parser := fastjson.Parser{}
	data, _ := parser.ParseBytes(s)
	stream := string(data.GetStringBytes("stream"))
	if len(stream) == 0 {
		return
	}
	symbol := strings.Split(stream, "@")[0]
	tdata := futurestradedata{}
	tdata.TradeAmount, _ = strconv.ParseFloat(string(data.GetStringBytes("data", "q")), 64)
	tdata.TradePrice, _ = strconv.ParseFloat(string(data.GetStringBytes("data", "p")), 64)
	tdata.TradeTime = data.GetInt64("data", "E") / 1000
	strdata, _ := json.Marshal(&tdata)
	redisconn := reptile.RedisPool.Get()
	redisconn.Do("publish", "reptile_futures_trade", fmt.Sprint(symbol, "@", string(strdata)))
	redisconn.Do("publish", "reptile_futures_price", fmt.Sprint(symbol, "@", tdata.TradePrice))
	redisconn.Close()
}
