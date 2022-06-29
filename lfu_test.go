package localcache

import (
	"testing"
	"time"
)

func TestLFUCache(t *testing.T) {
	c := NewLFUCache(2, 10, 1)
	if err := c.set("name", []byte("tiger"), time.Second); err != nil {
		t.Log(err)
	}
	if err := c.set("name", []byte("tiger"), time.Second); err != nil {
		t.Log(err)
	}
	if err := c.set("age", []byte("25"), time.Second); err != nil {
		t.Log(err)
	}

	t.Log("name=", string(c.get("name")))
	t.Log("age=", string(c.get("age")))
	time.Sleep(time.Second * 3)
	t.Log("name=", string(c.get("name")))
	t.Log("age=", string(c.get("age")))
}

func TestMaxSlotsForLFUCache(t *testing.T) {
	c := NewLFUCache(2, 100, 1)
	if err := c.set("name", []byte("tiger"), time.Second); err != nil {
		t.Log(err)
	}
	if err := c.set("age", []byte("25"), time.Second); err != nil {
		t.Log(err)
	}
	if err := c.set("name", []byte("tiger"), time.Second); err != nil {
		t.Log(err)
	}
	if err := c.set("address", []byte("beijing"), time.Second); err != nil {
		t.Log(err)
	}

	t.Log("name=", string(c.get("name")))
	t.Log("age=", string(c.get("age")))
	t.Log("address=", string(c.get("address")))
	time.Sleep(time.Second * 3)
	t.Log("name=", string(c.get("name")))
	t.Log("age=", string(c.get("age")))
	t.Log("address=", string(c.get("address")))
}

func TestMaxSlotsForLFUCacheWithDiffValue(t *testing.T) {
	c := NewLFUCache(2, 100, 1)
	if err := c.set("name", []byte("tiger"), time.Second); err != nil {
		t.Log(err)
	}
	if err := c.set("age", []byte("25"), time.Second); err != nil {
		t.Log(err)
	}
	if err := c.set("name", []byte("tiger2"), time.Second); err != nil {
		t.Log(err)
	}
	if err := c.set("address", []byte("beijing"), time.Second); err != nil {
		t.Log(err)
	}

	t.Log("name=", string(c.get("name")))
	t.Log("age=", string(c.get("age")))
	t.Log("address=", string(c.get("address")))
	time.Sleep(time.Second * 3)
	t.Log("name=", string(c.get("name")))
	t.Log("age=", string(c.get("age")))
	t.Log("address=", string(c.get("address")))
}

func TestMaxMomoryForLFUCache(t *testing.T) {
	c := NewLFUCache(10, 7, 1)
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
	time.Sleep(time.Second * 3)
	t.Log("name=", string(c.get("name")))
	t.Log("age=", string(c.get("age")))
	t.Log("address=", string(c.get("address")))
}
