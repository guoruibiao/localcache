package localcache

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestLocalCache(t *testing.T) {
	Cacher = NewDefaultLocalCache(10, 10, 1)
	Set("name", []byte("tiger"), time.Second)
	t.Log(GetString("name"))
	time.Sleep(time.Second * 3)
	t.Log(GetString("name"))
}

func TestMaxSlots(t *testing.T) {
	Cacher = NewDefaultLocalCache(10, 10, 1)
	for idx := 0; idx < 11; idx++ {
		if err := Set(fmt.Sprintf("index_%d", idx), []byte("tiger"), time.Second); err != nil {
			t.Log(err.Error())
		}
	}
	if err := Set("index_11", []byte("tiger"), time.Second*3); err != nil {
		t.Log(err.Error())
	}
	time.Sleep(time.Second * 2)
	for idx := 0; idx < 11; idx++ {
		t.Log(GetString(fmt.Sprintf("index_%d", idx)))
	}
	t.Log(GetString("index_11"))
}

func TestMaxMemory(t *testing.T) {
	Cacher = NewDefaultLocalCache(10, 10, 1)
	for idx := 0; idx < 11; idx++ {
		if err := Set(fmt.Sprintf("index_%d", idx), []byte("tiger"), time.Second); err != nil {
			t.Log(err.Error())
		}
	}
	time.Sleep(time.Second * 2)
	for idx := 0; idx < 11; idx++ {
		t.Log(GetString(fmt.Sprintf("index_%d", idx)))
	}
	t.Log(GetString("index_11"))
}

func TestBenchLocalCache(t *testing.T) {

	runtime.GOMAXPROCS(4)
	Cacher = NewDefaultLocalCache(200000, 100000, 1)
	for i := 0; i < 1000000; i++ {
		idx := i % 50000
		key := fmt.Sprintf("k%d", idx)
		val := fmt.Sprintf("v%d", idx)
		Set(key, []byte(val), 1*time.Second)
		time.Sleep(10 * time.Microsecond)
	}
}
