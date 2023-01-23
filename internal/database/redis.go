package database

import (
	"fmt"
	"strconv"

	"github.com/monnand/goredis"
	"github.com/seaung/gagako/pkg/utils"
)

func RedisUnauthorized(host string, port int) bool {
	if utils.IsOpenPort(host, port) {
		var client goredis.Client
		port := strconv.Itoa(port)
		client.Addr = fmt.Sprintf("%s:%s", host, port)
		err := client.Set("test", []byte("ISOK"))
		if err != nil {
			utils.New().Warnning(fmt.Sprintf("expect error : %v", err))
			return false
		}

		res, err := client.Get("test")
		if string(res) == "ISOK" {
			client.Set("test", []byte("test"))
			return true
		}
	}
	return false
}

func RedisUnauthorized2(host string) {
	if RedisUnauthorized(host, 6379) {
		utils.New().Info(fmt.Sprintf("Found Redis Unauthorized Vaulnerable !"))
	}
}
