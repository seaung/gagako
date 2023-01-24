package memcached

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/seaung/gagako/pkg/utils"
)

func MemcachedUnauthorized(target string) {
	client := memcache.New(fmt.Sprintf("%s:11211", target))

	err := client.Ping()
	if err != nil {
		utils.New().Warnning(fmt.Sprintf("Connect to memcache error : %v\n", err))
		return
	}

	err = client.Set(&memcache.Item{
		Key:        "test",
		Value:      []byte("ISOK"),
		Flags:      0,
		Expiration: 100,
	})

	if err != nil {
		utils.New().Warnning(fmt.Sprintf("expect error : %v", err))
		return
	}

	value, err := client.Get("test")

	if err != nil {
		utils.New().Warnning("NOT FOUND MemcachedUnauthorized\n")
		return
	}

	if string(value.Value) == "ISOK" {
		utils.New().Success("FOUND MemcachedUnauthorized")
	}
}
