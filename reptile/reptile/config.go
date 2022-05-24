package reptile

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)
var Symbols = []string{}
var MarketStreamUrl string
var FututesStreamUrl string
var RedisPool redis.Pool
var KlineIntervals = []string{"1m","3m","5m","15m","30m","2h","4h","6h","8h","12h","1d","3d","1w","1M"}
var KlineIntervalSeconds = map[string]int{}
type SubScribeData struct {
	Id     int      `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}
func ConfigInit(){
	KlineIntervalSeconds["1m"]  = 60
	KlineIntervalSeconds["3m"]  = 60 * 3
	KlineIntervalSeconds["5m"]  = 60 * 5
	KlineIntervalSeconds["15m"] = 60 * 15
	KlineIntervalSeconds["30m"] = 60 * 30
	KlineIntervalSeconds["2h"]  = 60 * 60 * 2
	KlineIntervalSeconds["4h"]  = 60 * 60 * 4
	KlineIntervalSeconds["6h"]  = 60 * 60 * 6
	KlineIntervalSeconds["8h"]  = 60 * 60 * 8
	KlineIntervalSeconds["12h"] = 60 * 60 * 12
	KlineIntervalSeconds["1d"]  = 60 * 60 * 24
	KlineIntervalSeconds["3d"]  = 60 * 60 * 24 * 3
	KlineIntervalSeconds["1w"]  = 60 * 60 * 24 * 7
	viper.AddConfigPath("./")
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err:=viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	RedisHost := viper.GetString("server.redis.host")
	RedisPort := viper.GetInt("server.redis.port")
	RedisPassword := viper.GetString("server.redis.password")
	RedisDb := viper.GetInt("server.redis.db")
	Symbols = viper.GetStringSlice("server.symbols")
	MarketStreamUrl = viper.GetString("server.marketurl")
	FututesStreamUrl = viper.GetString("server.futuresurl")
 	RedisPool= redis.Pool{
		MaxIdle:     10,
		MaxActive:   100,
		IdleTimeout: time.Duration(60) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", fmt.Sprint(RedisHost, ":", RedisPort),
				redis.DialPassword(RedisPassword),
				redis.DialDatabase(RedisDb),
			)
			if err != nil {
				panic(err)
			}
			return con, nil
		},
	}
}