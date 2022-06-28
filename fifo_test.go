package localcache

import (
	"testing"
	"time"
)

func TestFIFOCache(t *testing.T) {
	c := NewFIFOCache(10, 10, 1)
	if err := c.set("name", []byte("tiger"), time.Second); err != nil {
		t.Log(err)
	}
	t.Log("name=", string(c.get("name")))
	time.Sleep(time.Second * 2)
	t.Log("name=", string(c.get("name")))
}

func TestMaxSlotsForFIFOCache(t *testing.T) {
	c := NewFIFOCache(2, 100, 1)
	if err := c.set("name", []byte("tiger"), time.Second); err != nil {
		t.Log(err)
	}
	if err := c.set("age", []byte("25"), time.Second); err != nil {
		t.Log(err)
	}
	if err := c.set("address", []byte("beijing"), time.Second); err != nil {
		t.Log(err)
	}
	t.Log("name=", string(c.get("name")))
	t.Log("age=", string(c.get("age")))
	t.Log("address=", string(c.get("address")))
	time.Sleep(time.Second * 2)
	t.Log("name=", string(c.get("name")))
	t.Log("age=", string(c.get("age")))
	t.Log("address=", string(c.get("address")))
}

func TestMaxMemoryForFIFOCache(t *testing.T) {
	c := NewFIFOCache(10, 7, 1)
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
	time.Sleep(time.Second * 2)
	t.Log("name=", string(c.get("name")))
	t.Log("age=", string(c.get("age")))
	t.Log("address=", string(c.get("address")))
}
