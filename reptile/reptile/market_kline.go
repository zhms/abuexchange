package reptile

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/imroc/req"
	"github.com/spf13/viper"
	"github.com/valyala/fastjson"
)

func MarketKLine() {
	for i := 0; i < len(Symbols); i++ {
		go market_kline(Symbols[i])
	}
}
func market_kline(symbol string){
	defer recover()
	conn, _, err := websocket.DefaultDialer.Dial(MarketStreamUrl, nil)
	if err != nil {
		fmt.Println("market kline:",err)
		time.Sleep(time.Second * 1)
		go market_kline(symbol)
		return
	}
	submsg := SubScribeData{}
	submsg.Id = 1
	submsg.Method = "SUBSCRIBE"
	submsg.Params = []string{}
	redisconn := RedisPool.Get()
	for i := 0; i < len(KlineIntervals); i++ {
		interval := KlineIntervals[i]
		submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol,"/","",-1),fmt.Sprint("@kline_",interval)))
		redisconn.Do("del",fmt.Sprint("reptile:market:kline:",strings.Replace(symbol,"/","",-1),":",interval,":lastid"))
		go market_kline_history(symbol,interval)
	}
	redisconn.Close()
	sendmsg,_ := json.Marshal(&submsg)
	conn.WriteMessage(1, sendmsg)
	go market_kline_cut(symbol)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		go market_kline_message(message)
	}
	time.Sleep(time.Second * 1)
	go market_kline(symbol)
}


type marketklinedata struct{
	Id int64 `json:"id"`
	Open float64 `json:"open"`
	Close float64 `json:"close"`
	High float64 `json:"high"`
	Low float64 `json:"low"`
	Vol float64 `json:"vol"`
	Count int64 `json:"count"`
	Amount float64 `json:"amount"`
}

func market_kline_message(s []byte) {
	defer recover()
	parser := fastjson.Parser{}
	data,_ := parser.ParseBytes(s)
	stream := string(data.GetStringBytes("stream"))
	if len(stream) == 0 { return }
	symbol := strings.Split(stream,"@")[0]
	interval := string(data.GetStringBytes("data","k","i"))
	kdata := marketklinedata{}
	kdata.Id = data.GetInt64("data","k","t") / 1000
	kdata.Open,_ = strconv.ParseFloat(string(data.GetStringBytes("data","k","o")),64)
	kdata.Close,_ = strconv.ParseFloat(string(data.GetStringBytes("data","k","c")),64)
	kdata.High,_ = strconv.ParseFloat(string(data.GetStringBytes("data","k","h")),64)
	kdata.Low,_ = strconv.ParseFloat(string(data.GetStringBytes("data","k","l")),64)
	kdata.Vol,_ = strconv.ParseFloat(string(data.GetStringBytes("data","k","v")),64)
	kdata.Amount,_ = strconv.ParseFloat(string(data.GetStringBytes("data","k","q")),64)
	kdata.Count = data.GetInt64("data","k","n")
	strdata,_ := json.Marshal(&kdata)
	rediskey := fmt.Sprint("reptile:market:kline:",symbol,":",interval)
	redisconn := RedisPool.Get()
	redisconn.Do("hset",rediskey,kdata.Id,strdata)
	redisconn.Do("set",fmt.Sprint(rediskey,":lastid"),kdata.Id)
	redisconn.Do("publish","reptile_market_kline",fmt.Sprint(symbol,"@",interval,"@",string(strdata)))
	redisconn.Close()
}

func market_kline_history(symbol string,interval string){
	redisconn := RedisPool.Get()
	for {
		lastsync,_ := redisconn.Do("get",fmt.Sprint("reptile:market:kline:",strings.Replace(symbol,"/","",-1),":",interval,":lastid"))
		if lastsync != nil{
			break
		}
		time.Sleep(time.Second * 1)
	}
	kkeys,_ := redisconn.Do("hkeys",fmt.Sprint("reptile:market:kline:",strings.Replace(symbol,"/","",-1),":",interval))
	ikkeys := kkeys.([]interface{})
	skkeys := []string{}
	for i := 0; i < len(ikkeys); i++ {
		skkeys = append(skkeys, string(ikkeys[i].([]uint8)))
	}
	sort.Slice(skkeys, func(i, j int) bool {
		a,_ := strconv.ParseInt(skkeys[i],10,64)
		b,_ := strconv.ParseInt(skkeys[j],10,64)
		return a > b
	})
	var synccount int = 0
	if len(skkeys) < 50 {
		synccount = 1000
	} else {
		if interval != "1M"  {
			seconds := KlineIntervalSeconds[interval]
			for i := 0; i < len(skkeys) - 1; i++ {
				a,_ := strconv.ParseInt(skkeys[i],10,64)
				b,_ := strconv.ParseInt(skkeys[i + 1],10,64)
				if a - b != int64(seconds) {
					synccount = 1000
					break
				}
			}
		}else{
			synccount = 1000
		}
	}
	if synccount > 0 {
		cfgurls:= viper.GetStringSlice("server.marketklinegeturls")
		if len(cfgurls)  > 0 {
			rnum := rand.Intn(len(cfgurls))
			url := fmt.Sprint(cfgurls[rnum],"/api/v3/klines?symbol=",strings.ToUpper(strings.Replace(symbol,"/","",-1)),"&interval=",interval,"&limit=",synccount)
			resp,err := req.Get(url)
			if err != nil{
				fmt.Println(err)
				return
			}
			var jdata []interface{}
			resp.ToJSON(&jdata)
			for i, j := 0, len(jdata)-1; i < j; i, j = i+1, j-1 {
				jdata[i], jdata[j] = jdata[j], jdata[i]
			}
			fmt.Println("market",url,len(jdata))
			for i := 0; i < len(jdata); i++ {
				jkline := jdata[i].([]interface{})
				kdata := marketklinedata{}
				kdata.Id = int64(jkline[0].(float64) / 1000)
				kdata.Open,_ = strconv.ParseFloat(jkline[1].(string),64)
				kdata.High,_ = strconv.ParseFloat(jkline[2].(string),64)
				kdata.Low,_ = strconv.ParseFloat(jkline[3].(string),64)
				kdata.Close,_ = strconv.ParseFloat(jkline[4].(string),64)
				kdata.Vol,_ = strconv.ParseFloat(jkline[5].(string),64)
				kdata.Amount,_ = strconv.ParseFloat(jkline[7].(string),64)
				kdata.Count = int64(jkline[8].(float64))
				rediskey := fmt.Sprint("reptile:market:kline:",strings.Replace(symbol,"/","",-1),":",interval)
				strdata,_ := json.Marshal(&kdata)
				redisconn.Do("hset",rediskey,kdata.Id,strdata)
			}
		}
	}
	redisconn.Close()
}

func market_kline_cut(symbol string){
	for{
		redisconn := RedisPool.Get()
		for i := 0; i < len(KlineIntervals); i++ {
			interval := KlineIntervals[i]
			rediskey := fmt.Sprint("reptile:market:kline:",strings.Replace(symbol,"/","",-1),":",interval)
			rlen,_ := redisconn.Do("hlen",rediskey)
			keylen := rlen.(int64)
			maxlen := 14400 //14400   1000
			if keylen > int64(maxlen) {
				kkeys,_ := redisconn.Do("hkeys",rediskey)
				ikkeys := kkeys.([]interface{})
				skkeys := []string{}
				for i := 0; i < len(ikkeys); i++ {
					skkeys = append(skkeys, string(ikkeys[i].([]uint8)))
				}
				sort.Slice(skkeys, func(i, j int) bool {
					a,_ := strconv.ParseInt(skkeys[i],10,64)
					b,_ := strconv.ParseInt(skkeys[j],10,64)
					return a > b
				})
				for j := 0; j < len(skkeys); j++ {
					if j > maxlen {
						redisconn.Do("hdel",rediskey,skkeys[j])
					}
				}
			}
		}
		redisconn.Close()
		time.Sleep(time.Hour * 24)
	}
}