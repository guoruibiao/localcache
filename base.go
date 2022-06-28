package localcache

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"time"
)

const (
	LocalCacheTypeDefault = "default"
	LocalCacheTypeLRU = "lru"
)

type Cache interface {
	delete(key string)
	set(key string, val []byte, ttl time.Duration) error
	get(key string) []byte
	statistic()
}

var Cacher Cache

type Configuration struct {
	Name string `toml:"name"`
	MaxSlots int `toml:"max_slots"`
	MaxMemory int `toml:"max_memory"`
	CacheStrategy string `toml:"cache_strategy"`
	CleanInterval int `toml:"clean_interval"`
}


type localValue struct {
	keyname string // just for lru-cache
	value *bytes.Buffer
	ttl time.Time
}

func (lv *localValue) isExpired() bool {
	return !lv.ttl.IsZero() && time.Now().Unix() > lv.ttl.Unix()
}


func Init(configPath string) error {
	if _, err := os.Stat(configPath); err != nil {
		panic(err)
		return err
	}
	var config Configuration
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		panic(err)
		return err
	}
	fmt.Println("init localcache config by [" + config.Name + "]")
	// 从配置文件中选择要使用的 类型
	switch config.CacheStrategy {
	case LocalCacheTypeLRU:
		Cacher = NewLRUCache(config.MaxSlots, config.MaxMemory, config.CleanInterval)
	default:
		Cacher = NewDefaultLocalCache(config.MaxSlots, config.MaxMemory, config.CleanInterval)
	}

	return nil
}

func Statistics() {
	Cacher.statistic()
}

func Del(key string) {
	Cacher.delete(key)
}

func Set(key string, val []byte, ttl time.Duration) error {
	return Cacher.set(key, val, ttl)
}

func Get(key string) (val []byte) {
	return Cacher.get(key)
}

func GetString(key string) (val string) {

	return string(Cacher.get(key))
}

// ... many other methods to be implemented ...