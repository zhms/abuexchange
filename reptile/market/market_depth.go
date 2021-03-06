package market

import (
	"encoding/json"
	"exreptile/reptile"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/valyala/fastjson"
)

func MarketDepth() {
	for i := 0; i < len(reptile.Symbols); i++ {
		symbol := reptile.Symbols[i]
		redisconn := reptile.RedisPool.Get()
		redisconn.Do("del", fmt.Sprint("reptile:market:depth:", strings.Replace(symbol, "/", "", -1), ":asks"))
		redisconn.Do("del", fmt.Sprint("reptile:market:depth:", strings.Replace(symbol, "/", "", -1), ":bids"))
		redisconn.Close()
		go market_depth(symbol)
	}
}

func market_depth(symbol string) {
	defer recover()
	conn, _, err := websocket.DefaultDialer.Dial(reptile.MarketStreamUrl, nil)
	if err != nil {
		fmt.Println("market depth:", err)
		time.Sleep(time.Second * 1)
		go market_depth(symbol)
		return
	}

	submsg := reptile.SubScribeData{}
	submsg.Id = 1
	submsg.Method = "SUBSCRIBE"
	submsg.Params = []string{}
	submsg.Params = append(submsg.Params, fmt.Sprint(strings.Replace(symbol, "/", "", -1), "@depth"))
	sendmsg, _ := json.Marshal(&submsg)
	conn.WriteMessage(1, sendmsg)
	go market_depth_deal(symbol)
	go market_depth_delete(symbol)
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
	data, _ := parser.ParseBytes(s)
	if len(data.GetStringBytes("stream")) == 0 {
		return
	}
	stream := string(data.GetStringBytes("stream"))
	if len(stream) == 0 {
		return
	}
	symbol := strings.Split(stream, "@")[0]
	redisconn := reptile.RedisPool.Get()
	bids := data.GetArray("data", "b")
	for i := 0; i < len(bids); i++ {
		bdata := bids[i].GetArray()
		price := bdata[0]
		amount := bdata[1]
		rediskey := fmt.Sprint("reptile:market:depth:", symbol, ":bids")
		byteamount, _ := amount.StringBytes()
		stramount := string(byteamount)
		byteprice, _ := price.StringBytes()
		strprice := string(byteprice)
		famount, _ := strconv.ParseFloat(string(byteamount), 64)
		if famount > 0 {
			redisconn.Do("hset", rediskey, strprice, stramount)
		} else {
			redisconn.Do("hdel", rediskey, strprice)
		}
	}
	asks := data.GetArray("data", "a")
	for i := 0; i < len(asks); i++ {
		adata := asks[i].GetArray()
		price := adata[0]
		amount := adata[1]
		rediskey := fmt.Sprint("reptile:market:depth:", symbol, ":asks")
		byteamount, _ := amount.StringBytes()
		stramount := string(byteamount)
		byteprice, _ := price.StringBytes()
		strprice := string(byteprice)
		famount, _ := strconv.ParseFloat(string(byteamount), 64)
		if famount > 0 {
			redisconn.Do("hset", rediskey, strprice, stramount)
		} else {
			redisconn.Do("hdel", rediskey, strprice)
		}
	}
	redisconn.Close()
}

type marketdepthdealdata struct {
	Price  float64
	Amount float64
}

type marketdepthleveldata struct {
	Asks [][]float64 `json:"asks"`
	Bids [][]float64 `json:"bids"`
}

func market_depth_deal(symbol string) {
	for {
		time.Sleep(time.Second * 1)
		redisconn := reptile.RedisPool.Get()
		asks, _ := redisconn.Do("hgetall", fmt.Sprint("reptile:market:depth:", strings.Replace(symbol, "/", "", -1), ":asks"))
		bids, _ := redisconn.Do("hgetall", fmt.Sprint("reptile:market:depth:", strings.Replace(symbol, "/", "", -1), ":bids"))
		arrasks := asks.([]interface{})
		arrbids := bids.([]interface{})
		if len(arrasks) == 0 || len(arrbids) == 0 {
			continue
		}
		parrasks := []marketdepthdealdata{}
		for i := 0; i < len(arrasks); i += 2 {
			pd := marketdepthdealdata{}
			pd.Price, _ = strconv.ParseFloat(string(arrasks[i].([]uint8)), 64)
			pd.Amount, _ = strconv.ParseFloat(string(arrasks[i+1].([]uint8)), 64)
			parrasks = append(parrasks, pd)
		}
		sort.Slice(parrasks, func(i, j int) bool {
			return parrasks[i].Price < parrasks[j].Price
		})
		parrbids := []marketdepthdealdata{}
		for i := 0; i < len(arrbids); i += 2 {
			pd := marketdepthdealdata{}
			pd.Price, _ = strconv.ParseFloat(string(arrbids[i].([]uint8)), 64)
			pd.Amount, _ = strconv.ParseFloat(string(arrbids[i+1].([]uint8)), 64)
			parrbids = append(parrbids, pd)
		}
		sort.Slice(parrbids, func(i, j int) bool {
			return parrbids[i].Price > parrbids[j].Price
		})
		lev, _ := redisconn.Do("hget", "reptile:config:depth:market", strings.Replace(symbol, "/", "", -1))
		if lev == nil {
			continue
		}
		strlev := string(lev.([]uint8))
		splitlev := strings.Split(strlev, "@")
		for i := 0; i < len(splitlev); i++ {
			numlev, _ := strconv.ParseFloat(splitlev[i], 64)
			var dec int = 0
			splitnumlev := strings.Split(splitlev[i], ".")
			if len(splitnumlev) > 1 {
				dec = len(splitnumlev[1])
			}
			mapasks := make(map[float64]float64)
			for k := 0; k < len(parrasks); k++ {
				d := float64(numlev) * math.Pow(10, float64(dec))
				price := math.Floor(parrasks[k].Price * math.Pow(10, float64(dec)))
				price = math.Floor(math.Floor(price/d) * d)
				price = price / math.Pow(10, float64(dec))
				mapasks[price] += parrasks[k].Amount
			}
			mapbids := make(map[float64]float64)
			for k := 0; k < len(parrbids); k++ {
				d := float64(numlev) * math.Pow(10, float64(dec))
				price := math.Floor(parrbids[k].Price * math.Pow(10, float64(dec)))
				price = math.Floor(math.Floor(price/d) * d)
				price = price / math.Pow(10, float64(dec))
				mapbids[price] += parrbids[k].Amount
			}
			publishdata := marketdepthleveldata{}
			for k, v := range mapasks {
				publishdata.Asks = append(publishdata.Asks, []float64{k, v})
			}
			for k, v := range mapbids {
				publishdata.Bids = append(publishdata.Bids, []float64{k, v})
			}
			sort.Slice(publishdata.Asks, func(i, j int) bool {
				return publishdata.Asks[i][0] < publishdata.Asks[j][0]
			})
			sort.Slice(publishdata.Bids, func(i, j int) bool {
				return publishdata.Bids[i][0] > publishdata.Bids[j][0]
			})
			publishdatafinal := marketdepthleveldata{}
			for i := 0; i < len(publishdata.Asks); i++ {
				publishdatafinal.Asks = append(publishdatafinal.Asks, publishdata.Asks[i])
				if i >= 30 {
					break
				}
			}
			for i := 0; i < len(publishdata.Bids); i++ {
				publishdatafinal.Bids = append(publishdatafinal.Bids, publishdata.Bids[i])
				if i >= 30 {
					break
				}
			}
			strpublishdata, _ := json.Marshal(&publishdatafinal)
			rediskey := fmt.Sprint("reptile:market:depth:", strings.Replace(symbol, "/", "", -1), ":", numlev)
			redisconn.Do("set", rediskey, strpublishdata)
		}
		redisconn.Close()
	}
}

func market_depth_delete(symbol string) {
	for {
		time.Sleep(time.Second * 1)
		redisconn := reptile.RedisPool.Get()
		{
			askkeys, _ := redisconn.Do("hkeys", fmt.Sprint("reptile:market:depth:", strings.Replace(symbol, "/", "", -1), ":asks"))
			arraskkeys := askkeys.([]interface{})
			strasks := []string{}
			if len(arraskkeys) > 5000 {
				for i := range arraskkeys {
					strasks = append(strasks, string(arraskkeys[i].([]uint8)))
				}
			}
			sort.Slice(strasks, func(i, j int) bool {
				a, _ := strconv.ParseFloat(strasks[i], 64)
				b, _ := strconv.ParseFloat(strasks[j], 64)
				return a > b
			})
			for i := 0; i < len(strasks); i++ {
				if i > 5000 {
					redisconn.Do("hdel", fmt.Sprint("reptile:market:depth:", strings.Replace(symbol, "/", "", -1), ":asks"), fmt.Sprint(strasks[i]))
				}
			}
		}
		{
			askbids, _ := redisconn.Do("hkeys", fmt.Sprint("reptile:market:depth:", strings.Replace(symbol, "/", "", -1), ":bids"))
			arraskbids := askbids.([]interface{})
			strbids := []string{}
			if len(arraskbids) > 5000 {
				for i := range arraskbids {
					strbids = append(strbids, string(arraskbids[i].([]uint8)))
				}
			}
			sort.Slice(strbids, func(i, j int) bool {
				a, _ := strconv.ParseFloat(strbids[i], 64)
				b, _ := strconv.ParseFloat(strbids[j], 64)
				return a < b
			})
			for i := 0; i < len(strbids); i++ {
				if i > 5000 {
					redisconn.Do("hdel", fmt.Sprint("reptile:market:depth:", strings.Replace(symbol, "/", "", -1), ":bids"), fmt.Sprint(strbids[i]))
				}
			}
		}
		redisconn.Close()
	}
}
