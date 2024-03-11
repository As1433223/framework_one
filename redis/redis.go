package redis

import (
	"fmt"
	"github.com/As1433223/framework_one/config"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v2"
)

type RedisConf struct {
	Addr string `yaml:"Addr"`
	Port int    `yaml:"Port"`
	Db   int    `yaml:"Db"`
}

func InitRedis(servername string, RedisFunc func(Redis *redis.Client) error) error {
	conf, err := config.GetConfig(servername)
	if err != nil {
		return err
	}
	var (
		RedisConfig RedisConf
		data        map[interface{}]interface{}
	)
	yaml.Unmarshal([]byte(conf), &data)
	RedisConfig.Addr = data["RedisConf"].(map[interface{}]interface{})["Addr"].(string)
	RedisConfig.Port = data["RedisConf"].(map[interface{}]interface{})["Port"].(int)
	RedisConfig.Db = data["RedisConf"].(map[interface{}]interface{})["Db"].(int)
	Redis := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", RedisConfig.Addr, RedisConfig.Port),
		DB:   RedisConfig.Db,
	})
	defer Redis.Close()
	err = RedisFunc(Redis)
	if err != nil {
		return err
	}
	return nil
}
