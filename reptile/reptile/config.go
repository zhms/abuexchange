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
type SubScribeData struct {
	Id     int      `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}
func ConfigInit(){
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