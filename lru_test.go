package localcache

import (
	"testing"
	"time"
)

func TestLRUCache(t *testing.T) {
	Cacher = NewLRUCache(10, 10, 1)
	err := Cacher.set("name", []byte("tiger"), time.Second)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(Cacher.get("name")))
	time.Sleep(time.Second * 2)
	t.Log(string(Cacher.get("name")))
}

func TestLRUCache2(t *testing.T) {
	c := NewLRUCache(2, 100, 1)
	if err := c.set("name", []byte("tiger"), time.Second); err != nil {
		t.Log(err)
	}
	if err := c.set("age", []byte("25"), time.Second); err != nil {
		t.Log(err)
	}
	t.Log(string(c.get("name")))

	if err := c.set("address", []byte("beijing"), time.Second); err != nil {
		t.Log(err)
	}
	if err := c.set("school", []byte("pku"), time.Second); err != nil {
		t.Log(err)
	}
	c.statistic()

	t.Log("name=", string(c.get("name")))
	t.Log("age=", string(c.get("age")))
	t.Log("address=", string(c.get("address")))
	time.Sleep(time.Second * 2)
	t.Log("name=", string(c.get("name")))
	t.Log("school=", string(c.get("school")))
}

func TestLRUCacheMaxSlots(t *testing.T) {
	c := NewLRUCache(2, 1000, 1)
	if err := c.set("name", []byte("tiger"), time.Second); err != nil {
		t.Error(err)
		return
	}
	if err := c.set("age", []byte("25"), time.Second); err != nil {
		t.Error(err)
		return
	}
	t.Log("name=", string(c.get("name")))
	c.statistic()
	if err := c.set("address", []byte("beijing"), time.Second); err != nil {
		t.Error(err)
		return
	}
	c.statistic()
	t.Log("name=", string(c.get("name")))
	t.Log("address=", string(c.get("address")))
	time.Sleep(time.Second * 2)
	t.Log("address=", string(c.get("address")))
}

func TestLRUCacheMaxMemory(t *testing.T) {
	c := NewLRUCache(10, 10, 1)
	if err := c.set("name", []byte("tigertiger"), time.Second); err != nil {
		t.Log(err)
	}
	t.Log("name=", string(c.get("name")))
	if err := c.set("age", []byte("25"), time.Second); err != nil {
		t.Log(err)
	}
	t.Log("name=", string(c.get("name")))
	t.Log("age=", string(c.get("age")))
}

func TestLRUCacheShrink(t *testing.T) {
	c := NewLRUCache(10, 7, 1)
	if err := c.set("name", []byte("tiger"), time.Second); err != nil {
		t.Log(err)
	}
	if err := c.set("age", []byte("25"), time.Second); err != nil {
		t.Log(err)
	}
	if err := c.set("address", []byte("bj"), time.Second); err != nil {
		t.Log(err)
	}
	t.Log("name=", string(c.get("name")))
	t.Log("age=", string(c.get("age")))
	t.Log("address=", string(c.get("address")))
}

func TestLRUCacheShrink2(t *testing.T) {
	c := NewLRUCache(10, 7, 1)
	if err := c.set("name", []byte("tiger"), time.Second); err != nil {
		t.Log(err)
	}
	if err := c.set("age", []byte("25"), time.Second); err != nil {
		t.Log(err)
	}
	t.Log("name=", string(c.get("name")))
	if err := c.set("address", []byte("bj"), time.Second); err != nil {
		t.Log(err)
	}
	t.Log("name=", string(c.get("name")))
	t.Log("age=", string(c.get("age")))
	t.Log("address=", string(c.get("address")))
}

func TestLRUCacheBackgroundClean(t *testing.T) {
	c := NewLRUCache(10, 7, 1)
	if err := c.set("name", []byte("tiger"), time.Second); err != nil {
		t.Log(err)
	}
	t.Log("name=", string(c.get("name")))
	time.Sleep(time.Second * 3)
	t.Log("name=", string(c.get("name")))

}
