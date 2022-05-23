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

func MarketDepth() {
	for i := 0; i < len(Symbols); i++ {
		go market_depth(Symbols[i])
	}
}

func market_depth(symbol string) {
	defer recover()
	conn, _, err := websocket.DefaultDialer.Dial(MarketStreamUrl, nil)
	if err != nil {
		fmt.Println("market depth:", err)
		time.Sleep(time.Second * 1)
		go market_depth(symbol)
		return
	}
	redisconn := RedisPool.Get()
	redisconn.Do("del",fmt.Sprint("market:depth:",strings.Replace(symbol, "/", "", -1),":asks"))
	redisconn.Do("del",fmt.Sprint("market:depth:",strings.Replace(symbol, "/", "", -1),":bids"))
	redisconn.Close()
	submsg := SubScribeData{}
	submsg.Id = 1
	submsg.Method = "SUBSCRIBE"
	submsg.Params = []string{}
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol, "/", "", -1), "@depth"))
	sendmsg, _ := json.Marshal(&submsg)
	conn.WriteMessage(1, sendmsg)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		go market_depth_message(message)
	}
	time.Sleep(time.Second * 1)
	go market_depth(symbol)
}

func market_depth_message(s []byte) {
	defer recover()
	parser := fastjson.Parser{}
	data,_ := parser.ParseBytes(s)
	if len(data.GetStringBytes("stream")) == 0 { return }
	stream := string(data.GetStringBytes("stream"))
	if len(stream) == 0 { return }
	symbol := strings.Split(stream,"@")[0]
	redisconn := RedisPool.Get()
	bids := data.GetArray("data","b")
	for i := 0; i < len(bids); i++ {
		bdata := bids[i].GetArray()
		price := bdata[0]
		amount := bdata[1]
		rediskey := fmt.Sprint("reptile:market:depth:",symbol,":bids")
		byteamount,_ := amount.StringBytes()
		stramount := string(byteamount)
		byteprice,_ := price.StringBytes()
		strprice := string(byteprice)
		famount,_ := strconv.ParseFloat(string(byteamount),64)
		if famount > 0 {
			redisconn.Do("hset",rediskey,strprice,stramount)
		} else {
			redisconn.Do("hdel",rediskey,strprice)
		}
	}
	asks := data.GetArray("data","a")
	for i := 0; i < len(asks); i++ {
		adata := asks[i].GetArray()
		price := adata[0]
		amount := adata[1]
		rediskey := fmt.Sprint("reptile:market:depth:",symbol,":asks")
		byteamount,_ := amount.StringBytes()
		stramount := string(byteamount)
		byteprice,_ := price.StringBytes()
		strprice := string(byteprice)
		famount,_ := strconv.ParseFloat(string(byteamount),64)
		if famount > 0 {
			redisconn.Do("hset",rediskey,strprice,stramount)
		} else {
			redisconn.Do("hdel",rediskey,strprice)
		}
	}
	redisconn.Close()
}
